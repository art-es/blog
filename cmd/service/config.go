package main

import (
	"fmt"
	"os"
)

type Config struct {
	appEnv     string
	serviceUrl string
	jwtSecret  []byte
	pgConnect  string
	kafkaURL   string
}

func getConfig() *Config {
	appEnv := getenv("APP_ENV", "LOCAL")
	pgHost := getenv("PG_HOST", "127.0.0.1")
	pgPort := getenv("PG_PORT", "5432")
	pgUser := getenv("PG_USER", "postgres")
	pgPass := getenv("PG_USER", "postgres")
	pgDBName := getenv("PG_DBNAME", "postgres")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if appEnv == "PROD" {
			panic("JWT_SECRET cannot be empty on PROD")
		}
		jwtSecret = "secret"
	}

	return &Config{
		appEnv:     appEnv,
		serviceUrl: getenv("SERVICE_PORT", ":8080"),
		jwtSecret:  []byte(jwtSecret),
		pgConnect:  fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPass, pgHost, pgPort, pgDBName),
		kafkaURL:   getenv("KAFKA_URL", "127.0.0.1:9092"),
	}
}

func getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
