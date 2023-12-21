package cmd

import (
	"os"
	"time"

	"infra/common/db"

	"github.com/joho/godotenv"
)

const (
	HttpServerHost        = "HTTP_SERVER_HOST"
	HttpServerPort        = "HTTP_SERVER_PORT"
	DatabaseUri           = "DATABASE_URI"
	DatabaseDriver        = "DATABASE_DRIVER"
	DatabaseDialect       = "DATABASE_DIALECT"
	MaxOpenConnections    = "MAX_OPEN_CONNECTIONS"
	MaxIdleConnections    = "MAX_IDLE_CONNECTIONS"
	MaxConnectionLifetime = "MAX_CONNECTION_LIFETIME"
)

type HttpServerConfig struct {
	Host string
	Port string
}

type AppConfig struct {
	HttpServerConfig *HttpServerConfig
	DatabaseConfig   *db.DatabaseConfig
}

func BuildConfiguration() (*AppConfig, error) {

	_ = godotenv.Load(".env.local")
	_ = godotenv.Load()

	connectionLifetime, _ := time.ParseDuration("5m")

	return &AppConfig{
		&HttpServerConfig{
			envOrDefault(HttpServerHost, "localhost"),
			envOrDefault(HttpServerPort, "3000"),
		},
		&db.DatabaseConfig{
			Uri:                   os.Getenv(DatabaseUri),
			Driver:                os.Getenv(DatabaseDriver),
			Dialect:               os.Getenv(DatabaseDialect),
			MaxOpenConnections:    25,
			MaxIdleConnections:    25,
			MaxConnectionLifetime: connectionLifetime,
		},
	}, nil
}

func envOrDefault(name string, fallback string) string {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	return env
}
