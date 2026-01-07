package config

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/fikryfahrezy/forward/blog-api/internal/database"
	"github.com/fikryfahrezy/forward/blog-api/internal/http"
	"github.com/fikryfahrezy/forward/blog-api/internal/logger"
)

type Config struct {
	Server   http.Config
	Database database.Config
	Logger   logger.Config
}

func Load() Config {
	loadEnvFile(".env")

	return Config{
		Server: http.Config{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: database.Config{
			DSN: getEnv("DB_DSN", ""),
		},
		Logger: logger.Config{
			Level:  logger.ParseLevel(getEnv("LOG_LEVEL", "info")),
			Format: logger.ParseFormat(getEnv("LOG_FORMAT", "text")),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close .env file", slog.String("error", err.Error()))
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, "\"'")

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				slog.Error("Failed to set environment variable", slog.String("key", key), slog.String("error", err.Error()))
			}
		}
	}
}
