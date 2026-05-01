package odds

import "time"

// Sport identifies the sport/league category.
type Sport string

const (
	SportSoccer Sport = "soccer"
	SportNBA    Sport = "basketball_nba"
)

// MarketType classifies the betting market.
type MarketType string

const (
	MarketMoneyline MarketType = "moneyline"
	MarketSpread    MarketType = "spreads"
	MarketTotal     MarketType = "totals"
)

// Outcome identifies the side of a bet.
type Outcome string

const (
	OutcomeHome  Outcome = "home"
	OutcomeAway  Outcome = "away"
	OutcomeDraw  Outcome = "draw"
	OutcomeOver  Outcome = "over"
	OutcomeUnder Outcome = "under"
)

// Event represents a sports match.
type Event struct {
	ID        string    `json:"id"`
	Sport     Sport     `json:"sport"`
	League    string    `json:"league"`
	HomeTeam  string    `json:"homeTeam"`
	AwayTeam  string    `json:"awayTeam"`
	StartTime time.Time `json:"startTime"`
	InPlay    bool      `json:"inPlay"`
}

// MarketKey uniquely identifies a market within an event.
type MarketKey struct {
	EventID string     `json:"eventId"`
	Type    MarketType `json:"type"`
	Line    float64    `json:"line"`
}

// Price is a bookmaker's price for an outcome.
type Price struct {
	Book    string    `json:"book"`
	Decimal float64   `json:"decimal"`
	Updated time.Time `json:"updated"`
}

// Quote is a single bookmaker price for an outcome in a market.
type Quote struct {
	Key     MarketKey `json:"key"`
	Outcome Outcome   `json:"outcome"`
	Price   Price     `json:"price"`
}

// BestPrice is the best available price for an outcome across all bookmakers.
type BestPrice struct {
	Book    string  `json:"book"`
	Decimal float64 `json:"decimal"`
}
