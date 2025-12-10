package config

import (
	"os"
)

type Config struct {
	DBUrl     string
	Port      string
	JWTSecret string
	AppEnv    string
}

func Load() *Config {
	return &Config{
		DBUrl:     getEnv("DATABASE_URL", ""),
		Port:      getEnv("PORT", ":50051"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"), // Supabase JWT Secret
		AppEnv:    getEnv("APP_ENV", "production"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
