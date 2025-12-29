package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Addr:         getEnv("HTTP_ADDRESS", ":8080"),
			ReadTimeout:  getEnvDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
	}

	fmt.Println(cfg)
	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	fmt.Println(ok) // true because the env variable is injected using godotenv library

	if !ok {
		if fallback == "" {
			log.Fatalf("missing required env var: %s", key)
		}
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid duration value for env var %s: %s", key, value)
	}
	return duration
}
