package repository

import (
	"codegen/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Zero-allocation ipucu: Parametreleri pointer olarak al, gereksiz kopyalamadan kaçın.
func (r *PostgresRepo) CreateProduct(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	// pgx QueryRow performansı yüksektir.
	err := r.db.QueryRow(ctx, query, p.Name, p.Price, p.Stock).Scan(&p.ID)
	return err
}

func (r *PostgresRepo) GetProductByID(ctx context.Context, id int32) (*domain.Product, error) {
	p := &domain.Product{}
	query := `SELECT id, name, price, stock FROM products WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}
	return p, nil
}

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
