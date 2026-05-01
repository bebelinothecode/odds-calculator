# odds-calculator

OddsBot – a real-time odds comparison web app for NBA and major soccer leagues.

## Features

- Polls [The Odds API v4](https://the-odds-api.com) every 10–90 seconds for live and pre-match odds
- Compares prices across bookmakers for each event, market, and outcome
- Fires alerts when the best available price for an outcome improves by ≥ 2 %
- Streams alerts to the browser over SSE (Server-Sent Events)
- Serves a web UI with an odds comparison table and live alert feed

## Leagues covered

| Sport key | League |
|---|---|
| `basketball_nba` | NBA |
| `soccer_epl` | English Premier League |
| `soccer_uefa_champs_league` | UEFA Champions League |
| `soccer_spain_la_liga` | La Liga |
| `soccer_italy_serie_a` | Serie A |
| `soccer_germany_bundesliga` | Bundesliga |

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `ODDS_API_KEY` | ✅ | — | Your API key from [the-odds-api.com](https://the-odds-api.com) |
| `ODDS_API_BASE_URL` | ❌ | `https://api.the-odds-api.com` | Override the base URL (useful for testing) |
| `ODDS_API_SPORTKEYS` | ❌ | all 6 above | Comma-separated list of sport keys to poll |
| `SERVER_ADDR` | ❌ | `:8080` | TCP address the HTTP server listens on |

## Run

```bash
export ODDS_API_KEY="your_key_here"
go run ./cmd/oddsbot
```

Then open <http://localhost:8080>.

## Build

```bash
go build ./...
```

## API endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/events` | List current events (sorted by start time) |
| `GET` | `/api/events/{id}/odds` | Odds comparison for a specific event |
| `GET` | `/api/alerts/latest` | Latest best-price alerts (`?n=50` controls count) |
| `GET` | `/api/alerts/stream` | SSE stream of new best-price alerts |
| `GET` | `/` | Web UI |

## Project structure

```
cmd/oddsbot/           Entry point
internal/
  app/                 App wiring (New, Run) + polling loop
  config/              Environment-based configuration
  domain/
    odds/              Event, Quote, Price, MarketKey, BestPrice + compare helpers
    alerts/            Alert model
  providers/
    provider.go        Provider interface + EventWithOdds
    oddsapi/           The Odds API v4 adapter (client, mapper, provider)
  store/
    store.go           AlertStore + EventStore interfaces
    memory/            Thread-safe in-memory implementations
  transport/
    http/              HTTP routes, handlers, and SSE stream
web/                   Frontend (index.html)
```
