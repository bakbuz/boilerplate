package repository

import (
	"errors"
)

var ErrNotFound = errors.New("record not found")

const (
	failedToCollectRows = "failed to collect rows"
	listQueryRowError   = "list query row returned an error"
	rowScanError        = "row scan error"
	rowsIterationError  = "rows iteration error"
	failedToCount       = "failed to count records"
	failedToInsert      = "failed to insert record"
	failedToUpdate      = "failed to update record"
	failedToDelete      = "failed to delete record"
	failedToDeletes     = "failed to delete records"
	failedToSoftDelete  = "failed to soft delete record"
	failedToUpsert      = "failed to upsert record"
	failedToBulkInsert  = "failed to execute bulk insert"
	failedToBulkUpdate  = "failed to execute bulk update"
)

/*
// Transaction örneği
func (r *PostgresRepo) CreateOrder(ctx context.Context, o *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Stok düş
	_, err = tx.Exec(ctx, `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`, o.Quantity, o.ProductId)
	if err != nil {
		return err
	} // Stok yetersiz olabilir, bunu servis katmanı handle eder

	// 2. Sipariş oluştur
	query := `INSERT INTO orders (product_id, quantity) VALUES ($1, $2) RETURNING id, created_at`
	err = tx.QueryRow(ctx, query, o.ProductId, o.Quantity).Scan(&o.Id, &o.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type txContextType string

const txContextKey txContextType = "tx"

// Internal Helper: Context'te Tx varsa onu, yoksa Pool'u kullan
func (repo *productRepository) getDb(ctx context.Context) interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
} {
	if tx, ok := ctx.Value(txContextKey).(pgx.Tx); ok {
		return tx
	}
	return repo.db.Pool()
}
*/

const defaultBatchSize = 2000 // İdeal batch boyutu (2000-5000 arası genelde güvenlidir)
