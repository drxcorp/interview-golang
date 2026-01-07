package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	LogLevel   string
	AppEnv     string
}

var AppConfig *Config

func init() {
	AppConfig = &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "testdb"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		AppEnv:     getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetConfig() *Config {
	return AppConfig
}

func (c *Config) IsDevelopment() bool {
	if c.AppEnv == "development" {
		return true
	}
	return false
}

func (c *Config) IsProduction() bool {
	if c.AppEnv == "production" {
		return true
	}
	return false
}

func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName
}
