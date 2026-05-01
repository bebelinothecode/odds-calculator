package oddsapi

import "time"

// apiOddsEvent matches the /v4/sports/{sport_key}/odds response element.
type apiOddsEvent struct {
	ID           string    `json:"id"`
	SportKey     string    `json:"sport_key"`
	SportTitle   string    `json:"sport_title"`
	CommenceTime time.Time `json:"commence_time"`
	HomeTeam     string    `json:"home_team"`
	AwayTeam     string    `json:"away_team"`
	Bookmakers   []apiBook `json:"bookmakers"`
}

type apiBook struct {
	Key        string      `json:"key"`
	Title      string      `json:"title"`
	LastUpdate time.Time   `json:"last_update"`
	Markets    []apiMarket `json:"markets"`
}

type apiMarket struct {
	Key      string       `json:"key"` // h2h, spreads, totals
	Outcomes []apiOutcome `json:"outcomes"`
}

type apiOutcome struct {
	Name  string   `json:"name"`
	Price float64  `json:"price"` // decimal when oddsFormat=decimal
	Point *float64 `json:"point,omitempty"`
}
