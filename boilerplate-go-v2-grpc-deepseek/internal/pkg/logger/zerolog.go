package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
)

func New() zerolog.Logger {
    output := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
    }
    
    return zerolog.New(output).
        With().
        Timestamp().
        Caller().
        Logger().
        Level(zerolog.InfoLevel)
}

// Kullanım
log.Info().
    Str("method", "CreateProduct").
    Dur("duration", elapsed).
    Int("product_id", id).
    Msg("Product created successfully")