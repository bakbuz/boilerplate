package service

import (
	"codegen/internal/domain"
	"codegen/internal/repository"
	"codegen/pkg/errx"
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]*domain.Product, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Create(ctx context.Context, e *domain.Product) (int64, error)
	Update(ctx context.Context, e *domain.Product) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error)
	Count(ctx context.Context) (int64, error)
	BulkInsert(ctx context.Context, list []*domain.Product) (int64, error)
	Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error)
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService ...
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// validateProductId validates that product Id is not empty
func (s *productService) validateProductId(id uuid.UUID) error {
	if id == uuid.Nil {
		return errx.ErrInvalidInput
	}
	return nil
}

// GetAll retrieves all products
func (s *productService) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.GetAll(ctx)
}

// GetByIds retrieves products by their Ids
func (s *productService) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}

	return s.repo.GetByIds(ctx, ids)
}

// GetById retrieves a product by its Id
func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if err := s.validateProductId(id); err != nil {
		return nil, err
	}

	product, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errx.ErrNotFound
	}

	return product, nil
}

// Create creates a new product
func (s *productService) Create(ctx context.Context, e *domain.Product) (int64, error) {
	if err := e.Validate(); err != nil {
		return -1, err
	}

	// Generate new UUID if not provided
	if e.Id == uuid.Nil {
		e.Id, _ = uuid.NewV7()
	}

	return s.repo.Insert(ctx, e)
}

// Update updates an existing product
func (s *productService) Update(ctx context.Context, e *domain.Product) (int64, error) {
	if err := e.Validate(); err != nil {
		return -1, err
	}

	if err := s.validateProductId(e.Id); err != nil {
		return -1, err
	}

	return s.repo.Update(ctx, e)
}

// Delete permanently deletes a product
func (s *productService) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	if err := s.validateProductId(id); err != nil {
		return -1, err
	}

	return s.repo.Delete(ctx, id)
}

// SoftDelete soft deletes a product
func (s *productService) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	if err := s.validateProductId(id); err != nil {
		return -1, err
	}

	if deletedBy == uuid.Nil {
		return -1, errors.New("deletedBy is required for soft delete")
	}

	return s.repo.SoftDelete(ctx, id, deletedBy)
}

// Count returns total count of products
func (s *productService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// BulkInsert inserts multiple products
func (s *productService) BulkInsert(ctx context.Context, list []*domain.Product) (int64, error) {
	if len(list) == 0 {
		return -1, errx.ErrInvalidInput
	}

	// Validate all products before inserting
	for i, product := range list {
		if err := product.Validate(); err != nil {
			return -1, errors.Wrapf(err, "validation failed for product at index %d", i)
		}

		// Generate UUID if not provided
		if product.Id == uuid.Nil {
			product.Id, _ = uuid.NewV7()
		}
	}

	return s.repo.BulkInsert(ctx, list, 0)
}

// Search searches products based on filter criteria
func (s *productService) Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error) {
	if filter == nil {
		return nil, errx.ErrInvalidInput
	}

	// Validate pagination parameters
	if filter.Limit < 0 {
		return nil, errors.New("limit parameter must be non-negative")
	}

	// Set default pagination if not provided
	if filter.Limit == 0 {
		filter.Limit = 10 // Default page size
	}

	// Limit maximum page size to prevent excessive data retrieval
	if filter.Limit > 1000 {
		return nil, errors.New("limit parameter must not exceed 1000")
	}

	return s.repo.Search(ctx, filter)
}
