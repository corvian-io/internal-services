package managers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/corvian/argus/services/weather/internal/ingest"
)

type WeatherManagerI interface {
	Process(ctx context.Context) error
}

type WeatherManager struct {
	client ingest.Client
}

// Process implements [WeatherManagerI].
func (w *WeatherManager) Process(ctx context.Context) error {
	weather, err := w.client.Fetch(ctx)
	if err != nil {
		slog.Error("exception occured when processing weather data", "error", err)
		return err
	}
	fmt.Printf("got the weather data: %v\n", weather)
	
	return nil

}

func NewWeatherManager(client ingest.Client) WeatherManagerI {
	return &WeatherManager{
		client: client,
	}
}
