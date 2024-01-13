package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv            string
	ServiceURL        string
	AccessTokenSecret []byte
	PGConnect         string
	KafkaURL          string
}

func Parse() *Config {
	appEnv := getenv("APP_ENV", "LOCAL")
	pgHost := getenv("PG_HOST", "127.0.0.1")
	pgPort := getenv("PG_PORT", "5432")
	pgUser := getenv("PG_USER", "postgres")
	pgPass := getenv("PG_USER", "postgres")
	pgDBName := getenv("PG_DBNAME", "postgres")

	jwtSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if jwtSecret == "" {
		if appEnv == "PROD" {
			panic("ACCESS_TOKEN_SECRET cannot be empty on PROD")
		}
		jwtSecret = "secret"
	}

	return &Config{
		AppEnv:            appEnv,
		ServiceURL:        getenv("SERVICE_PORT", ":8080"),
		AccessTokenSecret: []byte(jwtSecret),
		PGConnect:         fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPass, pgHost, pgPort, pgDBName),
		KafkaURL:          getenv("KAFKA_URL", "127.0.0.1:9092"),
	}
}

func getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
