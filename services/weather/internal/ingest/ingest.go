package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/corvian/argus/services/weather/internal/config"
	model "github.com/corvian/argus/services/weather/internal/model"
)

type Client struct {
	config config.Config
}

func NewClient(config config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Fetch(ctx context.Context) ([]model.Weather, error) {
	//this is an example for now.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(c.config.GetCurrentUrl())
	if err != nil {
		slog.Error("exception occured returning weather service.", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("weather service didn't return a 200")
		return nil, fmt.Errorf("weather service didn't return a 200")
	}
	var w model.Weather
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		slog.Error("exception occured when decoding.", "error", err)
		return nil, err
	}
	return []model.Weather{w}, nil
}
