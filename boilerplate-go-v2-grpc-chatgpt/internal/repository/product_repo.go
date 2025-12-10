package repository

import (
	"codegen/internal/entity"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	Get(ctx context.Context, id int64) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id int64) error
}
