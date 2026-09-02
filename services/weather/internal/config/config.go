package config

import (
	"fmt"

	cfg "github.com/corvian/argus/internal/shared/configs"
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
	cfg.LoadEnvFileIfExists(".env")

	c := &Config{
		ApiKey:  cfg.GetEnvOrDefault(EnvKeys.ApiKey, ""),
		Url:     cfg.GetEnvOrDefault(EnvKeys.Url, ""),
		ZipCode: cfg.GetEnvOrDefault(EnvKeys.ZipCode, "21117"),
	}

	return c
}
