package configs

import "fmt"

type RedisConfig struct {
	Host     string
	Port     string
	Db       int
	Password string
}

func (r *RedisConfig) GetHostAddress() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
