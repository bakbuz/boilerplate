package configs

import (
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Log      LogConfig
}

type ServerConfig struct {
    Address                 string
    GRPCKeepAlive           time.Duration
    GRPCKeepAliveTime       time.Duration
    GRPCKeepAliveTimeout    time.Duration
    MaxMsgSize              int
    WorkerCount             int
    ShutdownTimeout         time.Duration
}

type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    SSLMode         string
    MaxConnections  int32
    MinConnections  int32
    MaxConnIdleTime time.Duration
    MaxConnLifetime time.Duration
    ConnectTimeout  time.Duration
}

type AuthConfig struct {
    JWTSecret     string
    TokenDuration time.Duration
}

type LogConfig struct {
    Level  string
    Format string
}

func Load() (*Config, error) {
    // Load from environment variables with defaults
    return &Config{
        Server: ServerConfig{
            Address:              ":50051",
            GRPCKeepAlive:        30 * time.Second,
            GRPCKeepAliveTime:    10 * time.Second,
            GRPCKeepAliveTimeout: 5 * time.Second,
            MaxMsgSize:           1024 * 1024 * 4, // 4MB
            WorkerCount:          10,
            ShutdownTimeout:      30 * time.Second,
        },
        Database: DatabaseConfig{
            Host:            "localhost",
            Port:            5432,
            User:            "postgres",
            Password:        "password",
            Database:        "products",
            SSLMode:         "require",
            MaxConnections:  10,
            MinConnections:  2,
            MaxConnIdleTime: 5 * time.Minute,
            MaxConnLifetime: 30 * time.Minute,
            ConnectTimeout:  10 * time.Second,
        },
        Auth: AuthConfig{
            JWTSecret:     "your-secret-key-change-in-production",
            TokenDuration: 24 * time.Hour,
        },
        Log: LogConfig{
            Level:  "info",
            Format: "json",
        },
    }, nil
}