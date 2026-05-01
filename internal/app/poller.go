package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"oddsbot/internal/domain/alerts"
	"oddsbot/internal/domain/odds"
	"oddsbot/internal/providers"
	"oddsbot/internal/store"
)

// bestKey identifies a unique (event, market, outcome) triple for tracking
// best-price changes between poll cycles.
type bestKey struct {
	EventID string
	Market  odds.MarketType
	Line    float64
	Outcome odds.Outcome
}

// alertSeq provides monotonically increasing alert IDs.
var alertSeq atomic.Uint64

// Poller continuously fetches events and quotes from a Provider, stores them,
// and fires alerts when the best price improves by at least minImprove.
type Poller struct {
	provider   providers.Provider
	alertStore store.AlertStore
	eventStore store.EventStore
	minImprove float64

	mu       sync.Mutex
	prevBest map[bestKey]odds.BestPrice
}

// NewPoller constructs a Poller. minImprove is a fractional threshold (e.g.
// 0.02 means ≥2% improvement triggers an alert).
func NewPoller(
	provider providers.Provider,
	alertStore store.AlertStore,
	eventStore store.EventStore,
	minImprove float64,
) *Poller {
	return &Poller{
		provider:   provider,
		alertStore: alertStore,
		eventStore: eventStore,
		minImprove: minImprove,
		prevBest:   make(map[bestKey]odds.BestPrice),
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
// Pre-match events are polled every 90 s; in-play every 10 s.
func (p *Poller) Run(ctx context.Context) {
	tPre := time.NewTicker(90 * time.Second)
	tLive := time.NewTicker(10 * time.Second)
	defer tPre.Stop()
	defer tLive.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tLive.C:
			p.poll(ctx)
		case <-tPre.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	now := time.Now()
	results, err := p.provider.FetchOdds(ctx, now.Add(-6*time.Hour), now.Add(24*time.Hour))
	if err != nil {
		log.Printf("poller: FetchOdds error: %v", err)
		return
	}

	events := make([]odds.Event, 0, len(results))
	for _, r := range results {
		events = append(events, r.Event)
		p.eventStore.SetQuotes(r.Event.ID, r.Quotes)
		p.checkAlerts(r.Event, r.Quotes)
	}
	if len(events) > 0 {
		p.eventStore.UpsertEvents(events)
	}
}

// checkAlerts computes the current best prices and fires alerts for outcomes
// whose best price improved by at least p.minImprove since the last poll.
func (p *Poller) checkAlerts(ev odds.Event, quotes []odds.Quote) {
	// Group quotes by market key.
	byMarket := map[odds.MarketKey][]odds.Quote{}
	for _, q := range quotes {
		byMarket[q.Key] = append(byMarket[q.Key], q)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for mk, qs := range byMarket {
		best := odds.BestByOutcome(qs)
		for outcome, bp := range best {
			k := bestKey{
				EventID: mk.EventID,
				Market:  mk.Type,
				Line:    mk.Line,
				Outcome: outcome,
			}

			old, seen := p.prevBest[k]
			p.prevBest[k] = bp

			if !seen || old.Decimal <= 0 {
				continue
			}
			if improve := (bp.Decimal - old.Decimal) / old.Decimal; improve >= p.minImprove {
				p.alertStore.Add(alerts.Alert{
					ID:        fmt.Sprintf("%d-%d", time.Now().UnixNano(), alertSeq.Add(1)),
					CreatedAt: time.Now(),
					EventID:   ev.ID,
					Market:    mk,
					Outcome:   outcome,
					Message:   "Best " + string(outcome) + " improved",
					OldBest:   old.Decimal,
					NewBest:   bp.Decimal,
					BestBook:  bp.Book,
				})
			}
		}
	}
}
