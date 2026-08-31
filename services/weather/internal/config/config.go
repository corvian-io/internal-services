package config

import (
	"fmt"
	"log/slog"
	"os"

	cfg "github.com/corvian/argus/internal/shared/configs"

	"github.com/lpernett/godotenv"
)

type Config struct {
	ApiKey  string
	Url     string
	ZipCode string
}

func (c *Config) GetCurrentUrl() string {
	return fmt.Sprintf("%s/current?access_key=%s&query=%s", c.Url, c.ApiKey, c.ZipCode)
}

func NewConfig() *Config {
	if _, err := os.Stat(".env"); err == nil {
		if err = godotenv.Load(".env"); err != nil {
			slog.Error("Error loading .env", "error", err)
		}
	}

	c := &Config{
		ApiKey:  cfg.GetEnvOrDefault(EnvKeys.ApiKey, ""),
		Url:     cfg.GetEnvOrDefault(EnvKeys.Url, ""),
		ZipCode: cfg.GetEnvOrDefault(EnvKeys.ZipCode, "21117"),
	}

	return c
}
