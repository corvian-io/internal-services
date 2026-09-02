package configs

import (
	"fmt"
	"strconv"

	"github.com/corvian/argus/internal/shared/constants"

	db "github.com/akhakpouri/gorm-kit/database"
)

type Config struct {
	Database db.DbConfig
	Redis    RedisConfig
}

func NewConfig() *Config {
	LoadEnvFileIfExists("configs/dev.env")
	portStr := GetEnvOrPanic(constants.EnvKeys.DBPort)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("invalid DB_PORT value: %s", portStr))
	}

	c := &Config{
		Database: db.DbConfig{
			Host:     GetEnvOrPanic(constants.EnvKeys.DBHost),
			Port:     port,
			User:     GetEnvOrPanic(constants.EnvKeys.DBUser),
			Password: GetEnvOrPanic(constants.EnvKeys.DBPassword),
			DbName:   GetEnvOrPanic(constants.EnvKeys.DBName),
			SSLMode:  GetEnvOrPanic(constants.EnvKeys.DBSSLMode),
			Schema:   GetEnvOrPanic(constants.EnvKeys.DBSchema),
		},
		Redis: RedisConfig{
			Host:     GetEnvOrPanic(constants.EnvKeys.RedisHost),
			Port:     GetEnvOrPanic(constants.EnvKeys.RedisPort),
			Db:       GetEnvOrDefaultToInt(constants.EnvKeys.RedisDb, 0),
			Password: GetEnvOrPanic(constants.EnvKeys.RedisPassword),
		},
	}

	return c
}

func NewDbConfig() *db.DbConfig {
	LoadEnvFileIfExists("configs/dev.env")

	portStr := GetEnvOrPanic(constants.EnvKeys.DBPort)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("invalid DB_PORT value: %s", portStr))
	}

	c := &db.DbConfig{
		Host:     GetEnvOrPanic(constants.EnvKeys.DBHost),
		Port:     port,
		User:     GetEnvOrPanic(constants.EnvKeys.DBUser),
		Password: GetEnvOrPanic(constants.EnvKeys.DBPassword),
		DbName:   GetEnvOrPanic(constants.EnvKeys.DBName),
		SSLMode:  GetEnvOrPanic(constants.EnvKeys.DBSSLMode),
		Schema:   GetEnvOrPanic(constants.EnvKeys.DBSchema),
	}

	return c
}

func NewRedisConfig() *RedisConfig {
	LoadEnvFileIfExists("configs/dev.env")
	c := RedisConfig{
		Host:     GetEnvOrPanic(constants.EnvKeys.RedisHost),
		Port:     GetEnvOrPanic(constants.EnvKeys.RedisPort),
		Db:       GetEnvOrDefaultToInt(constants.EnvKeys.RedisDb, 0),
		Password: GetEnvOrPanic(constants.EnvKeys.RedisPassword),
	}
	return &c
}
