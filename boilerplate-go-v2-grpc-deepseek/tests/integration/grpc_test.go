package integration_test

import (
    "context"
    "testing"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/ory/dockertest/v3"
    "github.com/stretchr/testify/assert"

    "github.com/yourusername/grpc-highperf-backend/internal/repository/postgres"
)

func TestProductRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    pool, err := dockertest.NewPool("")
    assert.NoError(t, err)
    
    resource, err := pool.Run("postgres", "15-alpine", []string{
        "POSTGRES_PASSWORD=secret",
        "POSTGRES_DB=testdb",
    })
    assert.NoError(t, err)
    defer pool.Purge(resource)
    
    var dbPool *pgxpool.Pool
    err = pool.Retry(func() error {
        connStr := "postgres://postgres:secret@localhost:" + resource.GetPort("5432/tcp") + "/testdb"
        cfg, err := pgxpool.ParseConfig(connStr)
        if err != nil {
            return err
        }
        
        dbPool, err = pgxpool.NewWithConfig(context.Background(), cfg)
        if err != nil {
            return err
        }
        
        return dbPool.Ping(context.Background())
    })
    assert.NoError(t, err)
    
    repo := postgres.NewProductRepository(dbPool)
    
    // Test cases...
}