package database

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog"
)

type PostgresConfig struct {
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

func NewPostgresPool(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
    connString := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
    )

    pgxCfg, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }

    // Performance optimizations for connection pooling
    pgxCfg.MaxConns = cfg.MaxConnections
    pgxCfg.MinConns = cfg.MinConnections
    pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime
    pgxCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
    pgxCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

    // Connection health checks
    pgxCfg.HealthCheckPeriod = 1 * time.Minute

    // Prepared statement caching
    pgxCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

    // Log connection info
    pgxCfg.ConnConfig.Tracer = &queryTracer{}

    pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
    if err != nil {
        return nil, err
    }

    // Test connection
    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }

    return pool, nil
}

// queryTracer for logging slow queries
type queryTracer struct{}

func (qt *queryTracer) TraceQueryStart(ctx context.Context, _ *pgconn.Conn, data pgx.TraceQueryStartData) context.Context {
    return context.WithValue(ctx, "query_start", time.Now())
}

func (qt *queryTracer) TraceQueryEnd(ctx context.Context, _ *pgconn.Conn, data pgx.TraceQueryEndData) {
    if start, ok := ctx.Value("query_start").(time.Time); ok {
        duration := time.Since(start)
        if duration > 100*time.Millisecond {
            logger := zerolog.Ctx(ctx)
            logger.Warn().
                Str("sql", data.SQL).
                Dur("duration", duration).
                Msg("Slow query detected")
        }
    }
}