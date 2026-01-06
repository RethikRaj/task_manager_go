package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTPConfig `env-prefix:"HTTP_"`
	DB   DBConfig
	Auth AuthConfig
}

// These tags replace the manual functions written before : getEnv, getEnvDuration , manual defaults , parsing logic
type HTTPConfig struct {
	Addr         string        `env:"ADDRESS" env-default:":8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

type AuthConfig struct {
	JWTSecret string `env:"JWT_SECRET" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	fmt.Println(cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	fmt.Println(cfg)
	return &cfg, nil
}
