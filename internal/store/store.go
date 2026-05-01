package store

import (
	"oddsbot/internal/domain/alerts"
	"oddsbot/internal/domain/odds"
)

// AlertStore persists and distributes best-price alerts.
type AlertStore interface {
	Add(a alerts.Alert)
	Latest(n int) []alerts.Alert
	// Subscribe returns a channel that receives new alerts and a cancel func
	// that unregisters the subscription and closes the channel.
	Subscribe() (<-chan alerts.Alert, func())
}

// EventStore persists events and their latest bookmaker quotes.
type EventStore interface {
	UpsertEvents(events []odds.Event)
	GetEvent(id string) (odds.Event, bool)
	ListEvents() []odds.Event
	SetQuotes(eventID string, quotes []odds.Quote)
	GetQuotes(eventID string) ([]odds.Quote, bool)
}
