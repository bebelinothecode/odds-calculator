package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"oddsbot/internal/domain/alerts"
	"oddsbot/internal/store"
)

// alertsLatest handles GET /api/alerts/latest
func alertsLatest(s store.AlertStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := 50
		if v := r.URL.Query().Get("n"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				n = parsed
			}
		}
		out := s.Latest(n)
		if out == nil {
			out = []alerts.Alert{}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}

// alertsStream handles GET /api/alerts/stream (Server-Sent Events).
func alertsStream(s store.AlertStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch, cancel := s.Subscribe()
		defer cancel()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case alert, ok := <-ch:
				if !ok {
					return
				}
				data, err := json.Marshal(alert)
				if err != nil {
					continue
				}
				fmt.Fprintf(w, "event: alert\ndata: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}
