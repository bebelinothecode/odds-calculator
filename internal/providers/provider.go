package providers

import (
	"context"
	"time"

	"oddsbot/internal/domain/odds"
)

// EventWithOdds bundles an event together with its current bookmaker quotes.
type EventWithOdds struct {
	Event  odds.Event
	Quotes []odds.Quote
}

// Provider fetches sports events and odds from an upstream data source.
type Provider interface {
	Name() string
	// FetchOdds returns all events and their current odds within the given
	// commence-time window.
	FetchOdds(ctx context.Context, since, until time.Time) ([]EventWithOdds, error)
}
