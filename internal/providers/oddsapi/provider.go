package oddsapi

import (
	"context"
	"net/url"
	"strings"
	"time"

	"oddsbot/internal/domain/odds"
	"oddsbot/internal/providers"
)

// Provider implements providers.Provider using The Odds API v4.
type Provider struct {
	client    *Client
	sportKeys []string
	regions   []string
	markets   []string
}

// Compile-time interface check.
var _ providers.Provider = (*Provider)(nil)

// NewProvider constructs a Provider that will poll the given sport keys.
// Regions default to us,uk,eu,au (broad coverage); markets default to
// h2h,spreads,totals.
func NewProvider(client *Client, sportKeys []string) *Provider {
	return &Provider{
		client:    client,
		sportKeys: sportKeys,
		regions:   []string{"us", "uk", "eu", "au"},
		markets:   []string{"h2h", "spreads", "totals"},
	}
}

func (p *Provider) Name() string { return "the-odds-api-v4" }

// FetchOdds calls the /v4/sports/{sport_key}/odds endpoint for every
// configured sport key and returns all events with their current quotes.
// since and until bound the commence_time filter.
func (p *Provider) FetchOdds(ctx context.Context, since, until time.Time) ([]providers.EventWithOdds, error) {
	var result []providers.EventWithOdds

	for _, sk := range p.sportKeys {
		raw, err := p.fetchOdds(ctx, sk, since, until)
		if err != nil {
			return nil, err
		}

		sport := sportFromKey(sk)
		for _, o := range raw {
			ev := odds.Event{
				ID:        o.ID,
				Sport:     sport,
				League:    o.SportTitle,
				HomeTeam:  o.HomeTeam,
				AwayTeam:  o.AwayTeam,
				StartTime: o.CommenceTime,
				InPlay:    isInPlay(o.CommenceTime),
			}
			result = append(result, providers.EventWithOdds{
				Event:  ev,
				Quotes: mapOddsToQuotes(o),
			})
		}
	}

	return result, nil
}

func (p *Provider) fetchOdds(ctx context.Context, sportKey string, since, until time.Time) ([]apiOddsEvent, error) {
	q := url.Values{}
	q.Set("regions", strings.Join(p.regions, ","))
	q.Set("markets", strings.Join(p.markets, ","))
	q.Set("oddsFormat", "decimal")
	q.Set("dateFormat", "iso")
	q.Set("commenceTimeFrom", since.UTC().Format(time.RFC3339))
	q.Set("commenceTimeTo", until.UTC().Format(time.RFC3339))

	var out []apiOddsEvent
	if err := p.client.doGET(ctx, "/v4/sports/"+sportKey+"/odds", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func mapOddsToQuotes(o apiOddsEvent) []odds.Quote {
	var out []odds.Quote
	for _, b := range o.Bookmakers {
		for _, m := range b.Markets {
			mt := marketFromKey(m.Key)
			for _, oc := range m.Outcomes {
				out = append(out, odds.Quote{
					Key: odds.MarketKey{
						EventID: o.ID,
						Type:    mt,
						Line:    asLine(oc.Point),
					},
					Outcome: outcomeFromName(oc.Name, o.HomeTeam, o.AwayTeam),
					Price: odds.Price{
						Book:    b.Key,
						Decimal: oc.Price,
						Updated: b.LastUpdate,
					},
				})
			}
		}
	}
	return out
}
