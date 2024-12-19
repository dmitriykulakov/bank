package config

import (
	"os"
)

type Config struct {
	ServerConfig
	DbConfig
}

type ServerConfig struct {
	Address string
}

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewConfig() *Config {
	return &Config{
		ServerConfig: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", "0.0.0.0:3351"),
		},
		DbConfig: DbConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Username: getEnv("POSTGRES_USERNAME", "kulakov"),
			Password: getEnv("POSTGRES_PASSWORD", "1234"),
			Database: getEnv("POSTGRES_NAME", "Bank"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
