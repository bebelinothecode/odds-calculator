package config

import (
	"fmt"
	"os"
	"strings"
)

// defaultSportKeys is the set of sport keys polled when ODDS_API_SPORTKEYS is
// not overridden.
var defaultSportKeys = []string{
	"basketball_nba",
	"soccer_epl",
	"soccer_uefa_champs_league",
	"soccer_spain_la_liga",
	"soccer_italy_serie_a",
	"soccer_germany_bundesliga",
}

// Config holds all runtime configuration.
type Config struct {
	OddsAPIKey     string
	OddsAPIBaseURL string
	SportKeys      []string
	ServerAddr     string
}

// Load reads configuration from environment variables.
// ODDS_API_KEY is required; all others have sensible defaults.
func Load() (*Config, error) {
	key := os.Getenv("ODDS_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("environment variable ODDS_API_KEY is required")
	}

	baseURL := os.Getenv("ODDS_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.the-odds-api.com"
	}

	sportKeys := defaultSportKeys
	if v := os.Getenv("ODDS_API_SPORTKEYS"); v != "" {
		sportKeys = splitComma(v)
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return &Config{
		OddsAPIKey:     key,
		OddsAPIBaseURL: baseURL,
		SportKeys:      sportKeys,
		ServerAddr:     addr,
	}, nil
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
