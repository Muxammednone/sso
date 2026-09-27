package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL     string        `envconfig:"DATABASE_URL" required:"true"`
	JWTSecret       string        `envconfig:"JWT_SECRET" required:"true"`
	GRPCAddr        string        `envconfig:"GRPC_ADDR" default:"127.0.0.1:50051"`
	TokenTTL        time.Duration `envconfig:"TOKEN_TTL" default:"1h"`
	DBTimeout       time.Duration `envconfig:"DB_TIMEOUT" default:"5s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

func Load() (Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	if config.GRPCAddr == "" {
		return Config{}, fmt.Errorf("GRPCAddr cant be empty")
	}
	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DatabaseURL cant be empty")
	}
	if len(config.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWTSecret string size must be minimum 32 bytes")
	}
	if config.TokenTTL <= 0 {
		return Config{}, fmt.Errorf("TokenTTL must be positive")
	}
	if config.DBTimeout <= 0 {
		return Config{}, fmt.Errorf("DBTimeout must be positive")
	}
	if config.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("ShutdownTimeout must be positive")
	}

	return config, nil
}
