package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	APIKey         string
	OpenAIAPIKey   string
	OpenAIModel    string
	EmbeddingModel string
	DatabasePath   string
	GinMode        string
}

var cfg *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: no .env file found, reading from environment")
	}

	cfg = &Config{
		Port:           getEnv("PORT", "8080"),
		APIKey:         getEnv("API_KEY", "dev-secret"),
		OpenAIAPIKey:   mustEnv("OPENAI_API_KEY"),
		OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		EmbeddingModel: getEnv("EMBEDDING_MODEL", "text-embedding-3-small"),
		DatabasePath:   getEnv("DATABASE_PATH", "data/coach.db"),
		GinMode:        getEnv("GIN_MODE", "debug"),
	}

	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func mustEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("FATAL: required env var %s is not set", key)
	}

	return value
}
