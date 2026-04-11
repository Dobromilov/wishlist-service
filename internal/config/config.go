package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

func Load() Config {
	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "db"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "postgres"),
		DBName:    getEnv("DB_NAME", "wishlist"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}
