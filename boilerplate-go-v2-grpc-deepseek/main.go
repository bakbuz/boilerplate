package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/yourusername/grpc-highperf-backend/internal/app"
    "github.com/yourusername/grpc-highperf-backend/internal/pkg/logger"
)

func main() {
    // Context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Setup logger
    log := logger.New()

    // Load configuration
    cfg, err := app.LoadConfig()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to load configuration")
    }

    // Create application
    application, err := app.NewApplication(ctx, cfg, log)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to create application")
    }

    // Start server
    go func() {
        if err := application.Start(ctx); err != nil {
            log.Error().Err(err).Msg("Server error")
        }
    }()

    // Wait for termination signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    // Graceful shutdown
    log.Info().Msg("Shutting down...")
    if err := application.Shutdown(ctx); err != nil {
        log.Error().Err(err).Msg("Error during shutdown")
    }
}