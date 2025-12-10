package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	// Config parse et
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	// t2.micro optimizasyonu: Bağlantı havuzunu sınırlı tut.
	// Çok fazla bağlantı RAM tüketir ve CPU context switch artırır.
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Connection Pool oluştur
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pool creation error: %w", err)
	}

	// Bağlantıyı test et (Ping)
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return pool, nil
}
