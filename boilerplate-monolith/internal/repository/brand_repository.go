package repository

import (
	"codegen/internal/database"
	"codegen/internal/entity"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// BrandRepository ...
type BrandRepository interface {
	GetAll(ctx context.Context) ([]*entity.Brand, error)
	GetByIds(ctx context.Context, ids []int32) ([]*entity.Brand, error)
	GetById(ctx context.Context, id int32) (*entity.Brand, error)
	Insert(ctx context.Context, e *entity.Brand) error
	Update(ctx context.Context, e *entity.Brand) (int64, error)
	Delete(ctx context.Context, id int32) (int64, error)
	Count(ctx context.Context) (int64, error)

	Upsert(ctx context.Context, e *entity.Brand) error
	BulkInsert(ctx context.Context, list []*entity.Brand) error
}

type brandRepository struct {
	db *database.DB
}

// NewBrandRepository ...
func NewBrandRepository(db *database.DB) BrandRepository {
	return &brandRepository{db: db}
}

func scanBrand(row pgx.Row) (*entity.Brand, error) {
	e := &entity.Brand{}

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
func (repo *brandRepository) GetAll(ctx context.Context) ([]*entity.Brand, error) {
	const stmt string = "SELECT * FROM catalog.brands"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*entity.Brand, error) {
		return scanBrand(row)
	})
	if err != nil {
		return nil, errors.WithMessage(err, failedToCollectRows)
	}

	return list, nil
}

// GetByIds ...
func (repo *brandRepository) GetByIds(ctx context.Context, ids []int32) ([]*entity.Brand, error) {
	if len(ids) == 0 {
		return []*entity.Brand{}, nil
	}

	const stmt string = "SELECT * FROM catalog.brands WHERE id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*entity.Brand, error) {
		return scanBrand(row)
	})
	if err != nil {
		return nil, errors.WithMessage(err, failedToCollectRows)
	}

	return list, nil
}

// GetById ...
func (repo *brandRepository) GetById(ctx context.Context, id int32) (*entity.Brand, error) {
	const stmt string = "SELECT * FROM catalog.brands WHERE id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := scanBrand(row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *brandRepository) Insert(ctx context.Context, e *entity.Brand) error {
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
func (repo *brandRepository) Update(ctx context.Context, e *entity.Brand) (int64, error) {
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

// Count ...
func (repo *brandRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM catalog.brands")
}

// Upsert ...
func (repo *brandRepository) Upsert(ctx context.Context, e *entity.Brand) error {
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
func (repo *brandRepository) BulkInsert(ctx context.Context, list []*entity.Brand) error {
	if len(list) == 0 {
		return nil
	}

	const command string = `
		INSERT INTO catalog.brands (name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	batch := &pgx.Batch{}
	for _, e := range list {
		batch.Queue(command,
			e.Name,
			e.Slug,
			e.Logo,
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
