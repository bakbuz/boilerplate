package repository

const (
	failedToCollectRows = "failed to collect rows"
	listQueryRowError   = "list query row returned an error"
	rowScanError        = "row scan error"
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
	_, err = tx.Exec(ctx, `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`, o.Quantity, o.ProductID)
	if err != nil {
		return err
	} // Stok yetersiz olabilir, bunu servis katmanı handle eder

	// 2. Sipariş oluştur
	query := `INSERT INTO orders (product_id, quantity) VALUES ($1, $2) RETURNING id, created_at`
	err = tx.QueryRow(ctx, query, o.ProductID, o.Quantity).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
*/
