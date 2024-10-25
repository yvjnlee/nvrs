package config

import (
	"log"
	"os"
)

var (
	JWTSecret []byte
	DB_URL    string
	Env       string
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Load() {
	DB_URL = getEnv("DB_URL", "sqlite3://./dev.db") // Default to SQLite for development
	Env = getEnv("ENV", "development")              // Default to "development" if ENV is not set

	log.Printf("Running in %s environment", Env)
	log.Printf("Database URL: %s", DB_URL)
}
