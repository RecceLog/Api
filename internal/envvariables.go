package internal

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err.Error())
	}
}

func GetString(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
