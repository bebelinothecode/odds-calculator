package httpapi

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"oddsbot/internal/domain/odds"
	"oddsbot/internal/store"
)

// eventsList handles GET /api/events
func eventsList(s store.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		evs := s.ListEvents()
		sort.Slice(evs, func(i, j int) bool {
			return evs[i].StartTime.Before(evs[j].StartTime)
		})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(evs)
	}
}

// eventOdds handles GET /api/events/{id}/odds
func eventOdds(s store.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expect path: /api/events/{id}/odds
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 4 || parts[0] != "api" || parts[1] != "events" || parts[3] != "odds" {
			http.NotFound(w, r)
			return
		}
		id := parts[2]

		ev, ok := s.GetEvent(id)
		if !ok {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}

		quotes, _ := s.GetQuotes(id)
		view := buildEventOddsView(ev, quotes)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(view)
	}
}

// ---- view model ----

type outcomeView struct {
	Outcome odds.Outcome       `json:"outcome"`
	Best    odds.BestPrice     `json:"best"`
	Prices  map[string]float64 `json:"prices"` // book -> decimal
}

type marketView struct {
	Type     odds.MarketType `json:"type"`
	Line     float64         `json:"line"`
	Outcomes []outcomeView   `json:"outcomes"`
}

type eventOddsView struct {
	Event   odds.Event   `json:"event"`
	Markets []marketView `json:"markets"`
	Books   []string     `json:"books"`
}

func buildEventOddsView(ev odds.Event, quotes []odds.Quote) eventOddsView {
	books := odds.SortedBooks(quotes)

	type mk struct {
		t odds.MarketType
		l float64
	}

	byMarket := map[mk][]odds.Quote{}
	for _, q := range quotes {
		key := mk{q.Key.Type, q.Key.Line}
		byMarket[key] = append(byMarket[key], q)
	}

	marketKeys := make([]mk, 0, len(byMarket))
	for k := range byMarket {
		marketKeys = append(marketKeys, k)
	}
	sort.Slice(marketKeys, func(i, j int) bool {
		if marketKeys[i].t == marketKeys[j].t {
			return marketKeys[i].l < marketKeys[j].l
		}
		return string(marketKeys[i].t) < string(marketKeys[j].t)
	})

	markets := make([]marketView, 0, len(marketKeys))
	for _, k := range marketKeys {
		qs := byMarket[k]
		best := odds.BestByOutcome(qs)

		outcomeSet := map[odds.Outcome]struct{}{}
		for _, q := range qs {
			outcomeSet[q.Outcome] = struct{}{}
		}
		outcomesList := make([]odds.Outcome, 0, len(outcomeSet))
		for o := range outcomeSet {
			outcomesList = append(outcomesList, o)
		}
		sort.Slice(outcomesList, func(i, j int) bool {
			return string(outcomesList[i]) < string(outcomesList[j])
		})

		outViews := make([]outcomeView, 0, len(outcomesList))
		for _, oc := range outcomesList {
			prices := map[string]float64{}
			for _, q := range qs {
				if q.Outcome == oc {
					prices[q.Price.Book] = q.Price.Decimal
				}
			}
			outViews = append(outViews, outcomeView{
				Outcome: oc,
				Best:    best[oc],
				Prices:  prices,
			})
		}

		markets = append(markets, marketView{
			Type:     k.t,
			Line:     k.l,
			Outcomes: outViews,
		})
	}

	return eventOddsView{
		Event:   ev,
		Markets: markets,
		Books:   books,
	}
}
