package repository

import (
	"codegen/internal/database"
	"codegen/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// ProductRepository ...
type ProductRepository interface {
	GetAll(ctx context.Context) ([]*domain.Product, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Insert(ctx context.Context, e *domain.Product) (int64, error)
	Update(ctx context.Context, e *domain.Product) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error)
	Count(ctx context.Context) (int64, error)

	Upsert(ctx context.Context, e *domain.Product) error
	BulkInsert(ctx context.Context, list []*domain.Product) (int64, error)
	BulkUpdate(ctx context.Context, list []*domain.Product) (int64, error)
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
	Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error)
}

type productRepository struct {
	db *database.DB
}

// NewProductRepository ...
func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db}
}

func scanProduct(row pgx.Row) (*domain.Product, error) {
	e := &domain.Product{}

	err := row.Scan(&e.Id, &e.BrandId, &e.Name, &e.Sku, &e.Summary, &e.Storyline, &e.StockQuantity, &e.Price, &e.Deleted, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt, &e.DeletedBy, &e.DeletedAt)
	if err == pgx.ErrNoRows { // sql: no rows in result set
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithMessage(err, rowScanError)
	}
	return e, nil
}

// GetAll ...
func (repo *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, product)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetByIds ...
func (repo *productRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}

	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false AND id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, product)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetById ...
func (repo *productRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false AND id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := scanProduct(row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *productRepository) Insert(ctx context.Context, e *domain.Product) (int64, error) {
	const command string = `
		INSERT INTO catalog.products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.Deleted,
		e.CreatedBy,
		e.CreatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToInsert)
	}

	return result.RowsAffected(), nil
}

// Update ...
func (repo *productRepository) Update(ctx context.Context, e *domain.Product) (int64, error) {
	const command string = `
		UPDATE catalog.products 
		SET brand_id=$2, name=$3, sku=$4, summary=$5, storyline=$6, stock_quantity=$7, price=$8, updated_by=$9, updated_at=$10 
		WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToUpdate)
	}

	return result.RowsAffected(), nil
}

// Delete ...
func (repo *productRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	const command string = `DELETE FROM catalog.products WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDelete)
	}

	return result.RowsAffected(), nil
}

// DeleteByIds ...
func (repo *productRepository) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	const command string = `DELETE FROM catalog.products WHERE id = ANY($1)`

	result, err := repo.db.Pool().Exec(ctx, command, ids)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDeletes)
	}

	return result.RowsAffected(), nil
}

// SoftDelete ...
func (repo *productRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	const command string = `UPDATE catalog.products SET deleted=true, deleted_by=$2, deleted_at=$3 WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id, deletedBy, time.Now().UTC())
	if err != nil {
		return -1, errors.WithMessage(err, failedToSoftDelete)
	}

	return result.RowsAffected(), nil
}

// Count ...
func (repo *productRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM catalog.products WHERE deleted=false")
}

// Upsert ...
func (repo *productRepository) Upsert(ctx context.Context, e *domain.Product) error {
	if e.Id == uuid.Nil {
		var err error
		if e.Id, err = uuid.NewV7(); err != nil {
			return errors.WithMessage(err, "failed to generate id")
		}
	}

	const command string = `
		INSERT INTO catalog.products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE 
		SET brand_id=EXCLUDED.brand_id,
			name=EXCLUDED.name,
			sku=EXCLUDED.sku,
			summary=EXCLUDED.summary,
			storyline=EXCLUDED.storyline,
			stock_quantity=EXCLUDED.stock_quantity,
			price=EXCLUDED.price,
			deleted=EXCLUDED.deleted,
			updated_by=$12,
			updated_at=$13`

	_, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.Deleted,
		e.CreatedBy,
		e.CreatedAt,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return errors.WithMessage(err, failedToUpsert)
	}

	return nil
}

func (repo *productRepository) BulkInsert(ctx context.Context, list []*domain.Product) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	rows := make([][]any, len(list))
	for i, e := range list {
		if e.Id == uuid.Nil {
			newId, err := uuid.NewV7()
			if err != nil {
				return 0, errors.WithMessage(err, "failed to generate id")
			}
			e.Id = newId
		}
		rows[i] = []any{e.Id, e.BrandId, e.Name, e.Sku, e.Summary, e.Storyline, e.StockQuantity, e.Price, e.Deleted, e.CreatedBy, e.CreatedAt}
	}

	count, err := repo.db.Pool().CopyFrom(
		ctx,
		pgx.Identifier{"catalog", "products"},
		[]string{"id", "brand_id", "name", "sku", "summary", "storyline", "stock_quantity", "price", "deleted", "created_by", "created_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkInsert)
	}

	return count, nil
}

func (repo *productRepository) BulkUpdate(ctx context.Context, list []*domain.Product) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Create temp table
	_, err = tx.Exec(ctx, `
		CREATE TEMP TABLE tmp_products_update (
			id uuid,
			brand_id int,
			name text,
			sku text,
			summary text,
			storyline text,
			stock_quantity int,
			price numeric,
			updated_by uuid,
			updated_at timestamp
		) ON COMMIT DROP
	`)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to create temp table")
	}

	// Prepare rows for CopyFrom
	rows := make([][]any, len(list))
	for i, e := range list {
		rows[i] = []any{
			e.Id,
			e.BrandId,
			e.Name,
			e.Sku,
			e.Summary,
			e.Storyline,
			e.StockQuantity,
			e.Price,
			e.UpdatedBy,
			e.UpdatedAt,
		}
	}

	// Copy data into temp table
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_products_update"},
		[]string{"id", "brand_id", "name", "sku", "summary", "storyline", "stock_quantity", "price", "updated_by", "updated_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to copy to temp table")
	}

	// Execute Update from temp table
	cmdTag, err := tx.Exec(ctx, `
		UPDATE catalog.products p
		SET 
			brand_id = t.brand_id,
			name = t.name,
			sku = t.sku,
			summary = t.summary,
			storyline = t.storyline,
			stock_quantity = t.stock_quantity,
			price = t.price,
			updated_by = t.updated_by,
			updated_at = t.updated_at
		FROM tmp_products_update t
		WHERE p.id = t.id
	`)
	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkUpdate)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return cmdTag.RowsAffected(), nil
}

func (repo *productRepository) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Context'e tx ekle
	txCtx := context.WithValue(ctx, txContextKey, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (repo *productRepository) Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error) {
	where, args := buildSearchWhere(filter)

	countQuery := "SELECT COUNT(*) FROM catalog.products WHERE deleted=false" + where

	// Get total count
	var total int64
	err := repo.db.Pool().QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, errors.WithMessage(err, failedToCount)
	}

	if total == 0 {
		return &domain.ProductSearchResult{}, nil
	}

	// Add pagination
	query := `SELECT p.id, p.name, p.price::numeric, p.brand_id, b.name AS brand_name
	FROM catalog.products AS p
	JOIN catalog.brands AS b ON p.brand_id = b.id
	WHERE p.deleted=false ` + where + ` 
	ORDER BY p.id DESC`

	argIndex := len(args) + 1
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	// Execute query
	rows, err := repo.db.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var items []domain.ProductSummary
	for rows.Next() {
		e := domain.ProductSummary{}
		err := rows.Scan(&e.Id, &e.Name, &e.Price, &e.BrandId, &e.BrandName)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		items = append(items, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return &domain.ProductSearchResult{
		Total: total,
		Items: items,
	}, nil
}

func buildSearchWhere(filter *domain.ProductSearchFilter) (string, []any) {
	var where string
	var args []any
	argIndex := 1

	if filter.BrandId > 0 {
		where += fmt.Sprintf(" AND brand_id=$%d", argIndex)
		args = append(args, filter.BrandId)
		argIndex++
	}

	if filter.Name != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}
	return where, args
}
