package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, connString string) (*DB, error) {
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
	config.ConnConfig.ConnectTimeout = 30 * time.Second

	// Connection health checks
	config.HealthCheckPeriod = 1 * time.Minute

	// Prepared statement caching
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Log connection info
	config.ConnConfig.Tracer = &queryTracer{}

	// Connection Pool oluştur
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pool creation error: %w", err)
	}

	// Bağlantıyı test et (Ping)
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Close
func (db *DB) Close() {
	db.pool.Close()
}

// Ping
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Pool'a erişim gerekirse
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Version(ctx context.Context) string {
	var version string
	err := db.pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return err.Error()
	}
	return version
}

func (db *DB) VersionWithErr(ctx context.Context) (string, error) {
	var version string
	err := db.pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return "", err
	}
	return version, nil
}

/**************************************************************************/
/***************************** BASE FUNCTIONS *****************************/
/**************************************************************************/

func GetDatabaseName(connString string) string {
	var dbName string
	var strs = strings.Split(connString, ";")
	for _, str := range strs {
		if strings.HasPrefix(strings.TrimSpace(str), "database") {
			dbName = strings.Split(str, "=")[1]
			break
		}
	}
	return dbName
}

func (db *DB) GetInt(ctx context.Context, stmt string, args ...any) (int, error) {
	var dest int

	row := db.pool.QueryRow(ctx, stmt, args...)
	if err := row.Scan(&dest); err != nil {
		if err == sql.ErrNoRows { // sql: no rows in result set
			return -1, errors.New("no rows")
		}
		return -1, err
	}
	return dest, nil
}

func (db *DB) GetString(ctx context.Context, stmt string, args ...any) (string, error) {
	var dest string

	row := db.pool.QueryRow(ctx, stmt, args...)
	if err := row.Scan(&dest); err != nil {
		if err == sql.ErrNoRows { // sql: no rows in result set
			return "", errors.New("no rows")
		}
		return "", err
	}

	return dest, nil
}

func (db *DB) Count(ctx context.Context, stmt string, args ...any) (int64, error) {
	var dest int64

	row := db.pool.QueryRow(ctx, stmt, args)
	if err := row.Scan(&dest); err != nil {
		if err == sql.ErrNoRows { // sql: no rows in result set
			return -1, errors.New("no rows")
		}
		return -1, err
	}

	return dest, nil
}

/**************************************************************************/
/****************************** QUERY TRACER ******************************/
/**************************************************************************/

// queryTracer for logging slow queries
type queryTracer struct{}

// TraceQueryStart implements pgx.QueryTracer.
func (qt *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, "query_start", time.Now())
}

// TraceQueryEnd implements pgx.QueryTracer.
func (qt *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if start, ok := ctx.Value("query_start").(time.Time); ok {
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			logger := zerolog.Ctx(ctx)
			logger.Warn().
				Str("sql", data.CommandTag.String()).
				Dur("duration", duration).
				Msg("Slow query detected")
		}
	}
}
