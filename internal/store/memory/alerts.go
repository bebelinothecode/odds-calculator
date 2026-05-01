package memory

import (
	"log"
	"sync"

	"oddsbot/internal/domain/alerts"
	"oddsbot/internal/store"
)

// Compile-time interface check.
var _ store.AlertStore = (*AlertStore)(nil)

// AlertStore is a thread-safe in-memory implementation of store.AlertStore.
type AlertStore struct {
	mu      sync.RWMutex
	alerts  []alerts.Alert
	subs    map[int]chan alerts.Alert
	nextSub int
}

// NewAlertStore returns an empty AlertStore.
func NewAlertStore() *AlertStore {
	return &AlertStore{
		subs: make(map[int]chan alerts.Alert),
	}
}

// Add appends an alert and fans it out to all subscribers.
func (s *AlertStore) Add(a alerts.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts = append(s.alerts, a)
	for _, ch := range s.subs {
		select {
		case ch <- a:
		default:
			log.Printf("alertstore: subscriber channel full, dropping alert id=%s", a.ID)
		}
	}
}

// Latest returns up to n of the most recent alerts (oldest first).
func (s *AlertStore) Latest(n int) []alerts.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.alerts) <= n {
		out := make([]alerts.Alert, len(s.alerts))
		copy(out, s.alerts)
		return out
	}
	src := s.alerts[len(s.alerts)-n:]
	out := make([]alerts.Alert, len(src))
	copy(out, src)
	return out
}

// Subscribe registers a new subscriber and returns a receive-only channel and
// a cancel function. Calling cancel unregisters the subscription and closes
// the channel.
func (s *AlertStore) Subscribe() (<-chan alerts.Alert, func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextSub
	s.nextSub++
	ch := make(chan alerts.Alert, 32)
	s.subs[id] = ch
	cancel := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.subs, id)
		close(ch)
	}
	return ch, cancel
}
