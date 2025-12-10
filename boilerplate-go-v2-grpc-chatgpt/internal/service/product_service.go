package service

import (
	"codegen/internal/entity"
	"codegen/internal/errorcodes"
	"codegen/internal/repository"
	"context"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) Create(ctx context.Context, p *entity.Product) error {
	// validation
	if p.Name == "" || p.Price <= 0 {
		return errorcodes.ErrInvalidInput
	}
	return s.repo.Create(ctx, p)
}
