package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/corvian/argus/services/weather/internal/config"
	"github.com/corvian/argus/services/weather/internal/ingest"
)

type Daemon struct {
	client ingest.Client
}

func NewDaemon(ctx context.Context, cfg *config.Config) (*Daemon, error) {
	return &Daemon{
		client: *ingest.NewClient(*cfg),
	}, nil
}

func (d *Daemon) Run(ctx context.Context) error {
	weather, err := d.client.Fetch(ctx)
	if err != nil {
		slog.Error("exception occured when processing weather data", "error", err)
		return err
	}
	fmt.Printf("got the weather data: %v\n", weather)
	return nil
}
