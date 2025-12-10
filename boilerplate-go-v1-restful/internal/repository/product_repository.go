package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"

	"codegen/internal/database"
	"codegen/internal/entity"
	"codegen/internal/repository/dto"
)

// ProductRepository ...
type ProductRepository interface {
	Search(ctx context.Context, filter *dto.ProductFilter) (int, []*entity.Product, error)
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Insert(ctx context.Context, e *entity.Product) (int64, error)
	Update(ctx context.Context, e *entity.Product) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy any) (int64, error)
	Count(ctx context.Context) (int64, error)

	BulkInsert(ctx context.Context, list []*entity.Product) error
}

type productRepository struct {
	db *database.DB
}

// NewProductRepository ...
func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db}
}

func mapToProducts(cursor *sql.Rows) ([]*entity.Product, error) {
	var entities []*entity.Product
	for cursor.Next() {
		e := &entity.Product{}

		err := cursor.Scan(&e.Id, &e.BrandId, &e.Name, &e.Sku, &e.Summary, &e.Storyline, &e.StockQuantity, &e.Price, &e.Deleted, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt, &e.DeletedBy, &e.DeletedAt)
		if err != nil {
			return nil, errors.WithMessage(err, "cursor scan error")
		}

		entities = append(entities, e)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.WithMessage(err, "cursor iteration error")
	}

	return entities, nil
}

func mapToProduct(row *sql.Row) (*entity.Product, error) {
	e := &entity.Product{}

	err := row.Scan(&e.Id, &e.BrandId, &e.Name, &e.Sku, &e.Summary, &e.Storyline, &e.StockQuantity, &e.Price, &e.Deleted, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt, &e.DeletedBy, &e.DeletedAt)
	if err == sql.ErrNoRows { // sql: no rows in result set
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithMessage(err, "row scan error")
	}
	return e, nil
}

// Search ...
func (repo *productRepository) Search(ctx context.Context, filter *dto.ProductFilter) (int, []*entity.Product, error) {
	const query = "SELECT * FROM Products WHERE Deleted=0"

	cursor, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return -1, nil, errors.WithMessage(err, "list query row returned an error")
	}
	defer cursor.Close()

	var list []*entity.Product

	total := 0
	if total > 0 {
		list, err = mapToProducts(cursor)
		if err != nil {
			return -1, nil, errors.WithStack(err)
		}
	}

	return total, list, nil
}

// GetAll ...
func (repo *productRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	const query = "SELECT * FROM Products WHERE Deleted=0"

	cursor, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.WithMessage(err, "list query row returned an error")
	}
	defer cursor.Close()

	list, err := mapToProducts(cursor)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

// GetByIds ...
func (repo *productRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	sids := fmt.Sprintf("'%s'", strings.Join(uuid.UUIDs.Strings(ids), "','"))
	query := fmt.Sprintf("SELECT * FROM Products WHERE Id IN (%s) AND Deleted=0", sids)

	cursor, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.WithMessage(err, "list query row returned an error")
	}
	defer cursor.Close()

	list, err := mapToProducts(cursor)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

// GetById ...
func (repo *productRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	const query = "SELECT * FROM Products WHERE Id=@Id AND Deleted=0"

	row := repo.db.QueryRowContext(ctx, query,
		sql.Named("Id", id),
	)
	if row.Err() != nil {
		return nil, errors.WithMessage(row.Err(), "single query row returned an error")
	}

	item, err := mapToProduct(row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *productRepository) Insert(ctx context.Context, e *entity.Product) (int64, error) {
	const query = `INSERT INTO Products (Id,BrandId,Name,Sku,Summary,Storyline,StockQuantity,Price,Deleted,CreatedBy,CreatedAt) VALUES (@Id,@BrandId,@Name,@Sku,@Summary,@Storyline,@StockQuantity,@Price,@Deleted,@CreatedBy,@CreatedAt)`

	result, err := repo.db.ExecContext(ctx, query,
		sql.Named("Id", e.Id),
		sql.Named("BrandId", e.BrandId),
		sql.Named("Name", e.Name),
		sql.Named("Sku", e.Sku),
		sql.Named("Summary", e.Summary),
		sql.Named("Storyline", e.Storyline),
		sql.Named("StockQuantity", e.StockQuantity),
		sql.Named("Price", e.Price),
		sql.Named("Deleted", e.Deleted),
		sql.Named("CreatedBy", e.CreatedBy),
		sql.Named("CreatedAt", e.CreatedAt),
	)
	if err != nil {
		return -1, errors.WithMessage(err, "insert command has returned an error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, errors.WithMessage(err, "insert command result number of rows affected returned error")
	}
	return rowsAffected, nil
}

// Update ...
func (repo *productRepository) Update(ctx context.Context, e *entity.Product) (int64, error) {
	const query = `UPDATE Products SET BrandId=@BrandId,Name=@Name,Sku=@Sku,Summary=@Summary,Storyline=@Storyline,StockQuantity=@StockQuantity,Price=@Price,UpdatedBy=@UpdatedBy,UpdatedAt=@UpdatedAt WHERE Id=@Id AND Deleted=0`

	result, err := repo.db.ExecContext(ctx, query,
		sql.Named("BrandId", e.BrandId),
		sql.Named("Name", e.Name),
		sql.Named("Sku", e.Sku),
		sql.Named("Summary", e.Summary),
		sql.Named("Storyline", e.Storyline),
		sql.Named("StockQuantity", e.StockQuantity),
		sql.Named("Price", e.Price),
		sql.Named("UpdatedBy", e.UpdatedBy),
		sql.Named("UpdatedAt", e.UpdatedAt),
		sql.Named("Id", e.Id),
	)
	if err != nil {
		return -1, errors.WithMessage(err, "update command has returned an error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, errors.WithMessage(err, "update command result number of rows affected returned error")
	}
	return rowsAffected, nil
}

// Delete ...
func (repo *productRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	const query = "DELETE FROM Products WHERE Id=@Id AND Deleted=0"

	result, err := repo.db.ExecContext(ctx, query,
		sql.Named("Id", id),
	)
	if err != nil {
		return -1, errors.WithMessage(err, "delete command has returned an error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, errors.WithMessage(err, "delete command result number of rows affected returned error")
	}
	return rowsAffected, nil
}

// SoftDelete ...
func (repo *productRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy any) (int64, error) {
	const query = "UPDATE Products SET Deleted=@Deleted,DeletedBy=@DeletedBy,DeletedAt=@DeletedAt WHERE Id=@Id"

	result, err := repo.db.ExecContext(ctx, query,
		sql.Named("Id", id),
		sql.Named("@DeletedBy", deletedBy),
		sql.Named("@DeletedAt", time.Now().UTC()),
	)
	if err != nil {
		return -1, errors.WithMessage(err, "update command has returned an error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, errors.WithMessage(err, "update command result number of rows affected returned error")
	}
	return rowsAffected, nil
}

// Count ...
func (repo *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Products WHERE Deleted=0")
	if err := row.Scan(&count); err != nil {
		return -1, errors.WithStack(err)
	}

	return count, nil
}

// BulkInsert ...
func (repo *productRepository) BulkInsert(ctx context.Context, list []*entity.Product) error {
	bulkOptions := mssql.BulkOptions{RowsPerBatch: 100}
	query := mssql.CopyIn("Products", bulkOptions, "Id", "BrandId", "Name", "Sku", "Summary", "Storyline", "StockQuantity", "Price", "Deleted", "CreatedBy", "CreatedAt")

	txn, err := repo.db.Begin()
	if err != nil {
		return errors.WithMessage(err, "begin wasn't starts a transaction.")
	}

	stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return errors.WithMessage(err, "txn wasn't prepare")
	}

	for _, e := range list {
		_, err := stmt.ExecContext(ctx, e.Id, e.BrandId, e.Name, e.Sku, e.Summary, e.Storyline, e.StockQuantity, e.Price, e.Deleted, e.CreatedBy, e.CreatedAt)
		if err != nil {
			return errors.WithMessage(err, "statement wasn't execute")
		}
	}

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		return errors.WithMessage(err, "statement execute")
	}

	if err = stmt.Close(); err != nil {
		return errors.WithMessage(err, "statement close")
	}

	if err = txn.Commit(); err != nil {
		return errors.WithMessage(err, "txn commit")
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "result rows affected")
	}

	if len(list) == int(rowCount) {
		log.Printf("%d row copied\n", rowCount)
	}

	return nil
}
