package oddsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is a thin wrapper around the The Odds API v4 HTTP API.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// NewClient constructs a Client from explicit credentials.
// baseURL should be "https://api.the-odds-api.com" in production.
func NewClient(baseURL, apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey must not be empty")
	}
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP:    &http.Client{Timeout: 20 * time.Second},
	}, nil
}

func (c *Client) doGET(ctx context.Context, path string, q url.Values, out any) error {
	if q == nil {
		q = url.Values{}
	}
	q.Set("apiKey", c.APIKey)

	u := c.BaseURL + path + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("oddsapi: HTTP %d for %s", res.StatusCode, u)
	}
	return json.NewDecoder(res.Body).Decode(out)
}
