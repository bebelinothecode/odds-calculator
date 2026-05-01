package memory

import (
	"sync"

	"oddsbot/internal/domain/odds"
	"oddsbot/internal/store"
)

// Compile-time interface check.
var _ store.EventStore = (*EventStore)(nil)

// EventStore is a thread-safe in-memory implementation of store.EventStore.
type EventStore struct {
	mu     sync.RWMutex
	events map[string]odds.Event
	quotes map[string][]odds.Quote
}

// NewEventStore returns an empty EventStore.
func NewEventStore() *EventStore {
	return &EventStore{
		events: make(map[string]odds.Event),
		quotes: make(map[string][]odds.Quote),
	}
}

// UpsertEvents stores or updates a batch of events.
func (s *EventStore) UpsertEvents(events []odds.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range events {
		s.events[e.ID] = e
	}
}

// GetEvent returns the event with the given ID.
func (s *EventStore) GetEvent(id string) (odds.Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	return e, ok
}

// ListEvents returns all stored events in arbitrary order.
func (s *EventStore) ListEvents() []odds.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]odds.Event, 0, len(s.events))
	for _, e := range s.events {
		out = append(out, e)
	}
	return out
}

// SetQuotes replaces the quotes for the given event.
func (s *EventStore) SetQuotes(eventID string, quotes []odds.Quote) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quotes[eventID] = quotes
}

// GetQuotes returns the latest quotes for the given event.
func (s *EventStore) GetQuotes(eventID string) ([]odds.Quote, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.quotes[eventID]
	return q, ok
}
