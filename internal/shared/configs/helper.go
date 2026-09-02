package configs

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}

	return value
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetEnvOrDefaultToInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		slog.Error("exception occured when converting string to int", "error", err)
		return defaultValue
	}

	return num
}

func LoadEnvFileIfExists(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		if err = godotenv.Load(filePath); err != nil {
			slog.Error(fmt.Sprintf("Error loading %s", filePath), "error", err)
		}
	}
}
