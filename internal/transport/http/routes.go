package httpapi

import (
	"net/http"

	"oddsbot/internal/store"
)

// Routes builds and returns the application's HTTP handler.
func Routes(alertStore store.AlertStore, eventStore store.EventStore) http.Handler {
	mux := http.NewServeMux()

	// Alerts
	mux.HandleFunc("/api/alerts/latest", alertsLatest(alertStore))
	mux.HandleFunc("/api/alerts/stream", alertsStream(alertStore))

	// Events + odds comparison
	mux.HandleFunc("/api/events", eventsList(eventStore))
	mux.HandleFunc("/api/events/", eventOdds(eventStore))

	// Serve the web UI
	mux.Handle("/", http.FileServer(http.Dir("web")))

	return mux
}
