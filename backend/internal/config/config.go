package config

import (
	"os"
)

type Config struct {
	Port     string
	DBDriver string
	DBDSN    string
}

// Load loads configuration from environment variables
func Load() Config {
	return Config{
		Port:     getEnv("PORT", "8080"),
		DBDriver: getEnv("DB_DRIVER", "mysql"),
		DBDSN:    getEnv("DB_DSN", "user:password@tcp(localhost:3306)/url_shortener"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
