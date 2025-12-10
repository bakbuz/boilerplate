package repository

import (
	"codegen/internal/domain"
	"context"
)

type Repository interface {
	CreateProduct(ctx context.Context, p *domain.Product) error
	GetProductByID(ctx context.Context, id int32) (*domain.Product, error)
	CreateOrder(ctx context.Context, o *domain.Order) error
}
