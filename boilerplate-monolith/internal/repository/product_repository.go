package repository

import (
	"codegen/internal/database"
	"codegen/internal/entity"
	"codegen/internal/repository/dto"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// ProductRepository ...
type ProductRepository interface {
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Insert(ctx context.Context, e *entity.Product) (int64, error)
	Update(ctx context.Context, e *entity.Product) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error)
	Count(ctx context.Context) (int64, error)

	BulkInsert(ctx context.Context, list []*entity.Product) error
	Search(ctx context.Context, filter *dto.ProductSearchFilter) (*dto.ProductSearchResult, error)
}

type productRepository struct {
	db *database.DB
}

// NewProductRepository ...
func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db}
}

func mapToProduct(row pgx.Row) (*entity.Product, error) {
	e := &entity.Product{}

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
func (repo *productRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	const stmt string = "SELECT * FROM products WHERE deleted=false"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[entity.Product])
	if err != nil {
		return nil, errors.WithMessage(err, failedToCollectRows)
	}

	return list, nil
}

// GetByIds ...
func (repo *productRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	if len(ids) == 0 {
		return []*entity.Product{}, nil
	}

	const stmt string = "SELECT * FROM products WHERE deleted=false AND id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[entity.Product])
	if err != nil {
		return nil, errors.WithMessage(err, failedToCollectRows)
	}

	return list, nil
}

// GetById ...
func (repo *productRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	const stmt string = "SELECT * FROM products WHERE deleted=false AND id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := mapToProduct(row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *productRepository) Insert(ctx context.Context, e *entity.Product) (int64, error) {
	const command string = `
		INSERT INTO products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
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
func (repo *productRepository) Update(ctx context.Context, e *entity.Product) (int64, error) {
	const command string = `
		UPDATE products 
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
	const command string = `DELETE FROM products WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDelete)
	}

	return result.RowsAffected(), nil
}

// SoftDelete ...
func (repo *productRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	const command string = `UPDATE products SET deleted=true, deleted_by=$2, deleted_at=$3 WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id, deletedBy, time.Now().UTC())
	if err != nil {
		return -1, errors.WithMessage(err, failedToSoftDelete)
	}

	return result.RowsAffected(), nil
}

// Count ...
func (repo *productRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM products WHERE deleted=false")
}

func (repo *productRepository) BulkInsert(ctx context.Context, list []*entity.Product) error {
	if len(list) == 0 {
		return nil
	}

	const command string = `
		INSERT INTO products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	batch := &pgx.Batch{}
	for _, e := range list {
		batch.Queue(command,
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
	}

	batchResults := repo.db.Pool().SendBatch(ctx, batch)
	defer batchResults.Close()

	for i := 0; i < len(list); i++ {
		_, err := batchResults.Exec()
		if err != nil {
			return errors.WithMessage(err, failedToBulkInsert)
		}
	}

	return nil
}

func (repo *productRepository) Search(ctx context.Context, filter *dto.ProductSearchFilter) (*dto.ProductSearchResult, error) {
	query := "SELECT * FROM products WHERE deleted=false"
	countQuery := "SELECT COUNT(*) FROM products WHERE deleted=false"
	args := []interface{}{}
	argIndex := 1

	// Build WHERE clauses dynamically based on filter
	if filter.Id != nil {
		query += fmt.Sprintf(" AND id=$%d", argIndex)
		countQuery += fmt.Sprintf(" AND id=$%d", argIndex)
		args = append(args, *filter.Id)
		argIndex++
	}

	if filter.BrandId > 0 {
		query += fmt.Sprintf(" AND brand_id=$%d", argIndex)
		countQuery += fmt.Sprintf(" AND brand_id=$%d", argIndex)
		args = append(args, filter.BrandId)
		argIndex++
	}

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	// Get total count
	var total int
	err := repo.db.Pool().QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, errors.WithMessage(err, failedToCount)
	}

	if total == 0 {
		return &dto.ProductSearchResult{
			Total: 0,
			Items: []*entity.Product{},
		}, nil
	}

	// Add pagination
	query += " ORDER BY created_at DESC"
	if filter.Take > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Take)
		argIndex++
	}

	if filter.Skip > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Skip)
	}

	// Execute query
	rows, err := repo.db.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[entity.Product])
	if err != nil {
		return nil, errors.WithMessage(err, failedToCollectRows)
	}

	return &dto.ProductSearchResult{
		Total: total,
		Items: items,
	}, nil
}
