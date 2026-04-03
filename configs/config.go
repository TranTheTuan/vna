// Package configs loads all application configuration from environment variables.
package configs

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	SwaggerEnabled bool       `env:"SWAGGER_ENABLED" envDefault:"true"`
	Database       Database   `envPrefix:"DATABASE_"`
	Auth           Auth       `envPrefix:"AUTH_"`
	ChatServer     ChatServer `envPrefix:"CHAT_SERVER_"`
}

// LoadConfig reads .env (if present) then environment variables and returns Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
