package config

import "os"

// Config представляет конфигурацию приложения
type Config struct {
	DatabaseURL string
	Port        string
}

// NewConfig создает новый экземпляр конфигурации с данными из окружения
func NewConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/postgres"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}
}
