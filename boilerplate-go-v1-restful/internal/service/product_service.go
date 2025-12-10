package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"codegen/internal/entity"
	"codegen/internal/repository"
	"codegen/internal/repository/dto"
)

// ProductService ...
type ProductService interface {
	SearchProducts(ctx context.Context, filter *dto.ProductFilter) (int, []*entity.Product, error)
	GetAllProducts(ctx context.Context) ([]*entity.Product, error)
	GetProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	CreateProduct(ctx context.Context, e *entity.Product) error
	UpdateProduct(ctx context.Context, e *entity.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	CountProduct(ctx context.Context) (int64, error)
}

type productService struct {
	repo        repository.ProductRepository
	currentUser CurrentUser
}

// NewProductService ...
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// SearchProducts ...
func (srv *productService) SearchProducts(ctx context.Context, filter *dto.ProductFilter) (int, []*entity.Product, error) {
	return srv.repo.Search(ctx, filter)
}

// GetAllProducts ...
func (srv *productService) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	return srv.repo.GetAll(ctx)
}

// GetProductsByIds ...
func (srv *productService) GetProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	if len(ids) == 0 {
		return nil, errors.New("product identifiers can't be zero")
	}
	return srv.repo.GetByIds(ctx, ids)
}

// GetProductById ...
func (srv *productService) GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	if id == uuid.Nil {
		return nil, errors.Errorf("invalid data value: %s", id)
	}

	return srv.repo.GetById(ctx, id)
}

// CreateProduct ...
func (srv *productService) CreateProduct(ctx context.Context, e *entity.Product) error {
	if e == nil {
		panic("product can't be null")
	}

	rowsAffected, err := srv.repo.Insert(ctx, e)
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		logger := zerolog.Ctx(ctx)
		logger.Info().Msgf("insert rows affected: %v", rowsAffected)
	}

	return nil
}

// UpdateProduct ...
func (srv *productService) UpdateProduct(ctx context.Context, e *entity.Product) error {
	if e == nil {
		panic("product can't be null")
	}

	rowsAffected, err := srv.repo.Update(ctx, e)
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		logger := zerolog.Ctx(ctx)
		logger.Info().Msgf("update rows affected: %v", rowsAffected)
	}
	return nil
}

// DeleteProduct ...
func (srv *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.Errorf("invalid data value: %s", id)
	}

	rowsAffected, err := srv.repo.SoftDelete(ctx, id, srv.currentUser.Id)
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		logger := zerolog.Ctx(ctx)
		logger.Info().Msgf("delete rows affected: %v", rowsAffected)
	}
	return nil
}

// CountProduct ...
func (srv *productService) CountProduct(ctx context.Context) (int64, error) {
	return srv.repo.Count(ctx)
}
