package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort string
	DBUrl   string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DBUrl: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_NAME", "subscriptions"),
			getEnv("DB_SSLMODE", "disable"),
		),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
