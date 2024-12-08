package config

import (
	"os"
)

type Config struct {
	ServerConfig
	PgConfig
}

type ServerConfig struct {
	Address string
}

type PgConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewConfig() *Config {
	return &Config{
		ServerConfig: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ""),
		},
		PgConfig: PgConfig{
			Host:     getEnv("POSTGRES_HOST", ""),
			Port:     getEnv("POSTGRES_PORT", ""),
			Username: getEnv("POSTGRES_USERNAME", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Database: getEnv("POSTGRES_NAME", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
