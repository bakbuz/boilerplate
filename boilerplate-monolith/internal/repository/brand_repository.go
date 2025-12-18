package repository

import (
	"codegen/internal/database"
	"codegen/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// BrandRepository ...
type BrandRepository interface {
	GetAll(ctx context.Context) ([]*domain.Brand, error)
	GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error)
	GetById(ctx context.Context, id int32) (*domain.Brand, error)
	Insert(ctx context.Context, e *domain.Brand) error
	Update(ctx context.Context, e *domain.Brand) (int64, error)
	Delete(ctx context.Context, id int32) (int64, error)
	DeleteByIds(ctx context.Context, ids []int32) (int64, error)
	Count(ctx context.Context) (int64, error)

	Upsert(ctx context.Context, e *domain.Brand) error
	BulkInsert(ctx context.Context, list []*domain.Brand) (int64, error)
	BulkUpdate(ctx context.Context, list []*domain.Brand) (int64, error)
	BulkInsertTran(ctx context.Context, list []*domain.Brand) error
	BulkUpdateTran(ctx context.Context, list []*domain.Brand) error
}

type brandRepository struct {
	db *database.DB
}

// NewBrandRepository ...
func NewBrandRepository(db *database.DB) BrandRepository {
	return &brandRepository{db: db}
}

func scanBrand(row pgx.Row) (*domain.Brand, error) {
	e := &domain.Brand{}

	err := row.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Or return specific ErrNotFound if preferred
		}
		return nil, errors.WithMessage(err, rowScanError)
	}
	return e, nil
}

// GetAll ...
func (repo *brandRepository) GetAll(ctx context.Context) ([]*domain.Brand, error) {
	// WARNING: Unbounded query. Added safety limit.
	// TODO: Update interface to support pagination.
	const stmt string = "SELECT * FROM catalog.brands LIMIT 1000"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Brand
	for rows.Next() {
		e := &domain.Brand{}
		err := rows.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetByIds ...
func (repo *brandRepository) GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error) {
	if len(ids) == 0 {
		return []*domain.Brand{}, nil
	}

	const stmt string = "SELECT * FROM catalog.brands WHERE id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Brand
	for rows.Next() {
		e := &domain.Brand{}
		err := rows.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetById ...
func (repo *brandRepository) GetById(ctx context.Context, id int32) (*domain.Brand, error) {
	const stmt string = "SELECT * FROM catalog.brands WHERE id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := scanBrand(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *brandRepository) Insert(ctx context.Context, e *domain.Brand) error {
	const command string = `
		INSERT INTO catalog.brands (name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := repo.db.Pool().QueryRow(ctx, command,
		e.Name,
		e.Slug,
		e.Logo,
		e.CreatedBy,
		e.CreatedAt,
	).Scan(&e.Id)

	if err != nil {
		return errors.WithMessage(err, failedToInsert)
	}

	return nil
}

// Update ...
func (repo *brandRepository) Update(ctx context.Context, e *domain.Brand) (int64, error) {
	const command string = `
		UPDATE catalog.brands 
		SET name=$2, slug=$3, logo=$4, updated_by=$5, updated_at=$6 
		WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.Name,
		e.Slug,
		e.Logo,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToUpdate)
	}

	return result.RowsAffected(), nil
}

// Delete ...
func (repo *brandRepository) Delete(ctx context.Context, id int32) (int64, error) {
	const command string = `DELETE FROM catalog.brands WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDelete)
	}

	return result.RowsAffected(), nil
}

// DeleteByIds ...
func (repo *brandRepository) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	const command string = `DELETE FROM catalog.brands WHERE id = ANY($1)`

	result, err := repo.db.Pool().Exec(ctx, command, ids)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDeletes)
	}

	return result.RowsAffected(), nil
}

// Count ...
func (repo *brandRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM catalog.brands")
}

// Upsert ...
func (repo *brandRepository) Upsert(ctx context.Context, e *domain.Brand) error {
	// If Id is missing, treating as Insert creates a record with ID 0 on some systems
	// or fails validation. For auto-increment, we must delegate to Insert.
	if e.Id == 0 {
		return repo.Insert(ctx, e)
	}

	const command string = `
		INSERT INTO catalog.brands (id, name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET name=EXCLUDED.name,
			slug=EXCLUDED.slug, 
			logo=EXCLUDED.logo, 
			updated_by=$7, 
			updated_at=$8`

	_, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.Name,
		e.Slug,
		e.Logo,
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

// BulkInsert ...
func (repo *brandRepository) BulkInsert(ctx context.Context, list []*domain.Brand) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	rows := make([][]any, len(list))
	for i, e := range list {
		rows[i] = []any{
			e.Name,
			e.Slug,
			e.Logo,
			e.CreatedBy,
			e.CreatedAt,
		}
	}

	count, err := repo.db.Pool().CopyFrom(
		ctx,
		pgx.Identifier{"catalog", "brands"},
		[]string{"name", "slug", "logo", "created_by", "created_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkInsert)
	}

	return count, nil
}

// BulkUpdate ...
func (repo *brandRepository) BulkUpdate(ctx context.Context, list []*domain.Brand) (int64, error) {
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
		CREATE TEMP TABLE tmp_brands_update (
			id int,
			name text,
			slug text,
			logo text,
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
			e.Name,
			e.Slug,
			e.Logo,
			e.UpdatedBy,
			e.UpdatedAt,
		}
	}

	// Copy data into temp table
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_brands_update"},
		[]string{"id", "name", "slug", "logo", "updated_by", "updated_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to copy to temp table")
	}

	// Execute Update from temp table
	cmdTag, err := tx.Exec(ctx, `
		UPDATE catalog.brands b
		SET 
			name = t.name,
			slug = t.slug,
			logo = t.logo,
			updated_by = t.updated_by,
			updated_at = t.updated_at
		FROM tmp_brands_update t
		WHERE b.id = t.id
	`)
	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkUpdate)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return cmdTag.RowsAffected(), nil
}

// Default batch size to prevent memory pressure
const batchSize = 1000

// BulkInsert ...
func (repo *brandRepository) BulkInsertTran(ctx context.Context, list []*domain.Brand) error {
	if len(list) == 0 {
		return nil
	}

	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const command string = `
		INSERT INTO catalog.brands (name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	// Chunking logic
	for i := 0; i < len(list); i += batchSize {
		batch := &pgx.Batch{}
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		for _, e := range list[i:end] {
			batch.Queue(command,
				e.Name,
				e.Slug,
				e.Logo,
				e.CreatedBy,
				e.CreatedAt,
			)
		}

		batchResults := tx.SendBatch(ctx, batch)

		for j := 0; j < (end - i); j++ {
			_, err := batchResults.Exec()
			if err != nil {
				batchResults.Close()
				return errors.WithMessage(err, failedToBulkInsert)
			}
		}

		if err := batchResults.Close(); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// BulkUpdate ...
func (repo *brandRepository) BulkUpdateTran(ctx context.Context, list []*domain.Brand) error {
	if len(list) == 0 {
		return nil
	}

	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const command string = `
		UPDATE catalog.brands 
		SET name=$2, slug=$3, logo=$4, updated_by=$5, updated_at=$6 
		WHERE id=$1`

	// Chunking logic
	for i := 0; i < len(list); i += batchSize {
		batch := &pgx.Batch{}
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		for _, e := range list[i:end] {
			batch.Queue(command,
				e.Id,
				e.Name,
				e.Slug,
				e.Logo,
				e.UpdatedBy,
				e.UpdatedAt,
			)
		}

		batchResults := tx.SendBatch(ctx, batch)

		for j := 0; j < (end - i); j++ {
			_, err := batchResults.Exec()
			if err != nil {
				batchResults.Close()
				return errors.WithMessage(err, failedToBulkUpdate)
			}
		}

		if err := batchResults.Close(); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
