package src

import (
	"os"
	"time"

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

type DatabaseConfig struct {
	Uri                   string
	Driver                string
	Dialect               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
}

type AppConfig struct {
	HttpServerConfig *HttpServerConfig
	DatabaseConfig   *DatabaseConfig
}

func BuildConfiguration() (*AppConfig, error) {
	// if len(os.Getenv("ENV_CONFIG_MAP")) > 0 {
	// 	_ = godotenv.Load(os.Getenv("ENV_CONFIG_MAP"))
	// }

	// env := envOrDefault("ENVIRONMENT", "local")
	// if env == "" {
	// 	env = "local"
	// }

	_ = godotenv.Load(".env.local")
	_ = godotenv.Load()

	// openConnections := os.Getenv(MaxOpenConnections) //, err := strconv.Atoi(envOrDefault(MaxOpenConnections, "25"))
	// // if err != nil {
	// // 	return nil, fmt.Errorf("MAX_OPEN_CONNECTIONS environment variable not set")
	// // }

	// idleConnections, err := strconv.Atoi(envOrDefault(MaxIdleConnections, "25"))
	// if err != nil {
	// 	return nil, fmt.Errorf("MAX_IDLE_CONNECTIONS environment variable not set")
	// }

	connectionLifetime, _ := time.ParseDuration("5m")

	return &AppConfig{
		&HttpServerConfig{
			envOrDefault(HttpServerHost, "localhost"),
			envOrDefault(HttpServerPort, "3000"),
		},
		&DatabaseConfig{
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
