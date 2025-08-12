package config

import "os"

type Config struct {
	DatabaseURL    string
	KarenAIToken   string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:  getEnvWithDefault("DATABASE_URL", "postgresql://root@localhost:26257/stockdb?sslmode=disable"),
		KarenAIToken: getEnvWithDefault("KAREN_AI_TOKEN", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MzQsImVtYWlsIjoiYXJpZWxmZXJtaW40MDJAZ21haWwuY29tIiwiZXhwIjoxNzUzNDU1NzYyLCJpZCI6IiIsInBhc3N3b3JkIjoiYHB3YC8qKi9GUk9NLyoqL3VzZXJzLyoqLy0tV0hFUkUvKiovdXNlcm5hbWUvKiovTElLRS8qKi8nJWFyaWVsZmVybWluNDAyQGdtYWlsLmNvbSUnLyoqLy0tIn0.JGtwIPv62kuaIxD1akO-dYSZUFqCsB8--yztFwG7Dew"),
		Port:         getEnvWithDefault("PORT", "8080"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}