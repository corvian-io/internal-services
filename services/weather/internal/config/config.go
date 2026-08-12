package config

import "fmt"

type Config struct {
	ApiKey string
	Url    string
}

func (c *Config) GetCurrentUrl(query string) string {
	return fmt.Sprintf("%s/current?access_key=%s&query=%s", c.Url, c.ApiKey, query)
}
