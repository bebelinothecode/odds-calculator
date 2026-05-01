package oddsapi

import (
	"strings"
	"time"

	"oddsbot/internal/domain/odds"
)

func sportFromKey(k string) odds.Sport {
	k = strings.ToLower(k)
	switch {
	case strings.Contains(k, "basketball") || strings.Contains(k, "nba"):
		return odds.SportNBA
	case strings.Contains(k, "soccer"):
		return odds.SportSoccer
	default:
		return odds.Sport(k)
	}
}

func marketFromKey(k string) odds.MarketType {
	switch k {
	case "h2h":
		return odds.MarketMoneyline
	case "spreads":
		return odds.MarketSpread
	case "totals":
		return odds.MarketTotal
	default:
		return odds.MarketType(k)
	}
}

func outcomeFromName(name, home, away string) odds.Outcome {
	n := strings.TrimSpace(strings.ToLower(name))
	h := strings.ToLower(home)
	a := strings.ToLower(away)

	switch n {
	case h:
		return odds.OutcomeHome
	case a:
		return odds.OutcomeAway
	case "draw", "tie", "x":
		return odds.OutcomeDraw
	case "over":
		return odds.OutcomeOver
	case "under":
		return odds.OutcomeUnder
	default:
		if strings.Contains(n, h) {
			return odds.OutcomeHome
		}
		if strings.Contains(n, a) {
			return odds.OutcomeAway
		}
		return odds.Outcome(n)
	}
}

func asLine(pt *float64) float64 {
	if pt == nil {
		return 0
	}
	return *pt
}

// isInPlay returns true if the event has already started (heuristic for v1).
func isInPlay(commence time.Time) bool {
	return time.Now().After(commence)
}
