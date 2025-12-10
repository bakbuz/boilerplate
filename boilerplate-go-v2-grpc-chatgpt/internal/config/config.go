package config

import "time"

type Config struct {
	Port            string
	DBUrl           string
	DBMaxConns      int32
	DBMinConns      int32
	JWTSecret       string
	SupabaseURL     string
	GracefulTimeout time.Duration
}

func LoadFromEnv() Config {
	return Config{
		Port:            getEnv("PORT", "50051"),
		DBUrl:           getEnv("DATABASE_URL", ""),
		DBMaxConns:      getEnvInt32("DB_MAX_CONNS", 10),
		DBMinConns:      getEnvInt32("DB_MIN_CONNS", 1),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		GracefulTimeout: 10 * time.Second,
	}
}
