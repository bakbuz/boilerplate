package pg

import (
	"codegen/internal/entity"
	"codegen/internal/errorcodes"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) Create(ctx context.Context, p *entity.Product) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO products (sku, name, description, price, stock) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`,
		p.SKU, p.Name, p.Description, p.Price, p.Stock)
	return row.Scan(&p.ID, &p.CreatedAt)
}

func (r *ProductRepo) Get(ctx context.Context, id int64) (*entity.Product, error) {
	p := &entity.Product{}
	row := r.pool.QueryRow(ctx, `SELECT id, sku, name, description, price, stock, created_at FROM products WHERE id=$1`, id)
	if err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
		if errors.Is(err, pgxpool.ErrNoRows) {
			return nil, errorcodes.ErrNotFound
		}
		return nil, err
	}
	return p, nil
}
