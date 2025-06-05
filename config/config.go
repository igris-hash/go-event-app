package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	ServerMode  string
	ServerPort  string
	LogLevel    string
	Database    DatabaseConfig
	JWTSecret   string
	CorsOrigins []string
}

// DatabaseConfig holds all database related configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// Load returns a new Config struct
func Load() *Config {
	return &Config{
		ServerMode: getEnv("SERVER_MODE", "debug"),
		ServerPort: getEnv("SERVER_PORT", "8000"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "db"),
			Port:     getEnv("DB_PORT", ""),
			Name:     getEnv("DB_NAME", "events.db"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret:   getEnv("JWT_SECRET", "your-256-bit-secret"),
		CorsOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"*"}),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return []string{value}
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolVal, err := strconv.ParseBool(value)
		if err == nil {
			return boolVal
		}
	}
	return defaultVal
}
