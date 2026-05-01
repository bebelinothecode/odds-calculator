package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"oddsbot/internal/config"
	"oddsbot/internal/providers/oddsapi"
	"oddsbot/internal/store/memory"
	httpapi "oddsbot/internal/transport/http"
)

// App holds the wired-up application components.
type App struct {
	cfg    *config.Config
	server *http.Server
	poller *Poller
}

// New wires up all components: Odds API client, in-memory stores, poller, and
// HTTP server.
func New(cfg *config.Config) (*App, error) {
	client, err := oddsapi.NewClient(cfg.OddsAPIBaseURL, cfg.OddsAPIKey)
	if err != nil {
		return nil, fmt.Errorf("creating odds api client: %w", err)
	}

	provider := oddsapi.NewProvider(client, cfg.SportKeys)

	alertStore := memory.NewAlertStore()
	eventStore := memory.NewEventStore()

	poller := NewPoller(provider, alertStore, eventStore, 0.02)

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           httpapi.Routes(alertStore, eventStore),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: srv,
		poller: poller,
	}, nil
}

// Run starts the poller in a goroutine and blocks on the HTTP server.
func (a *App) Run(ctx context.Context) error {
	go a.poller.Run(ctx)
	log.Printf("listening on http://localhost%s", a.cfg.ServerAddr)
	return a.server.ListenAndServe()
}
