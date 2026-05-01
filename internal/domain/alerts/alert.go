package alerts

import (
	"time"

	"oddsbot/internal/domain/odds"
)

// Alert fires when the best available price for an outcome improves by at
// least the configured threshold.
type Alert struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	EventID   string         `json:"eventId"`
	Market    odds.MarketKey `json:"market"`
	Outcome   odds.Outcome   `json:"outcome"`
	Message   string         `json:"message"`
	OldBest   float64        `json:"oldBest"`
	NewBest   float64        `json:"newBest"`
	BestBook  string         `json:"bestBook"`
}
