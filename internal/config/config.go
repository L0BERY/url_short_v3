package config

import (
	"os"
	"strings"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	DSN           string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress: normalizeAddress(getEnv("SERVER_ADDRESS", ":8080")),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		DSN:           getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/url_shortener_v3?sslmode=disable"),
	}
}

func normalizeAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ":8080"
	}
	if !strings.Contains(addr, ":") {
		return ":" + addr
	}
	return addr
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
