package config

import (
	"fmt"
	"log/slog"
	"os"

	cfg "github.com/corvian/argus/internal/shared/configs"

	"github.com/lpernett/godotenv"
)

type Config struct {
	ApiKey string
	Url    string
}

func (c *Config) GetCurrentUrl(query string) string {
	return fmt.Sprintf("%s/current?access_key=%s&query=%s", c.Url, c.ApiKey, query)
}

func NewConfig() *Config {
	if _, err := os.Stat("configs/dev.env"); err == nil {
		if err = godotenv.Load("configs/dev.env"); err != nil {
			slog.Error("Error loading configs/dev.env", "error", err)
		}
	}

	c := &Config{
		ApiKey: cfg.GetEnvOrDefault(EnvKeys.ApiKey, ""),
		Url:    cfg.GetEnvOrDefault(EnvKeys.Url, ""),
	}

	return c
}
