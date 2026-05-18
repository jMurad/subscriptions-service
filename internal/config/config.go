package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
}

func MustLoad() *Config {
	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
	}

	return cfg
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func mustEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}

	return value
}
