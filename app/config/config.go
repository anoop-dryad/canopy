package config

import (
	"os"
)

type Config struct {
	App  App
	HTTP HTTP
	DB   DB
}

type App struct {
	IsProduction bool
	Name         string
	Version      string
}

type HTTP struct {
	Addr string
}

type DB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

func Load() Config {
	return Config{
		App: App{
			IsProduction: getEnv("ENV", "development") == "production",
			Name:         getEnv("APP_NAME", "bridgehead"),
			Version:      getEnv("APP_VERSION", "dev"),
		},
		HTTP: HTTP{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		DB: DB{
			DSN:          mustEnv("DB_DSN"),
			MaxOpenConns: 25,
			MaxIdleConns: 5,
		},
	}
}

// mustEnv panics at startup if a required variable is missing.
// better to crash immediately than fail silently in production.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("required environment variable not set: " + key)
	}
	return val
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
