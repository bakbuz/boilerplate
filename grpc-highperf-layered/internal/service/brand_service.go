package service

import (
	"codegen/internal/domain"
	"codegen/internal/repository"
	"codegen/pkg/errx"
	"context"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type BrandService interface {
	GetAll(ctx context.Context) ([]*domain.Brand, error)
	GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error)
	GetById(ctx context.Context, id int32) (*domain.Brand, error)
	Create(ctx context.Context, e *domain.Brand) error
	Update(ctx context.Context, e *domain.Brand) (int64, error)
	Delete(ctx context.Context, id int32) (int64, error)
	Count(ctx context.Context) (int64, error)
	BulkInsertAll(ctx context.Context, list []*domain.Brand) (int64, error)
	BulkInsert(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error)
	BulkUpdate(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error)
}

type brandService struct {
	repo repository.BrandRepository
}

// NewBrandService ...
func NewBrandService(repo repository.BrandRepository) BrandService {
	return &brandService{repo: repo}
}

// validateBrand ...
func (s *brandService) validateBrand(e *domain.Brand) error {
	if e == nil {
		return errx.ErrInvalidInput
	}

	if strings.TrimSpace(e.Name) == "" {
		return errors.New("brand name is required")
	}

	if utf8.RuneCountInString(e.Name) > 255 {
		return errors.New("brand name must not exceed 255 characters")
	}

	if len(e.Slug) > 255 {
		return errors.New("brand slug must not exceed 255 characters")
	}

	return nil
}

// GetAll ...
func (s *brandService) GetAll(ctx context.Context) ([]*domain.Brand, error) {
	return s.repo.GetAll(ctx)
}

// GetByIds ...
func (s *brandService) GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error) {
	if len(ids) == 0 {
		return []*domain.Brand{}, nil
	}

	return s.repo.GetByIds(ctx, ids)
}

// GetById ...
func (s *brandService) GetById(ctx context.Context, id int32) (*domain.Brand, error) {
	if id == 0 {
		return nil, errx.ErrInvalidInput
	}

	brand, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if brand == nil {
		return nil, errx.ErrNotFound
	}

	return brand, nil
}

// Create ...
func (s *brandService) Create(ctx context.Context, e *domain.Brand) error {
	if err := s.validateBrand(e); err != nil {
		return err
	}

	return s.repo.Insert(ctx, e)
}

// Update ...
func (s *brandService) Update(ctx context.Context, e *domain.Brand) (int64, error) {
	if err := s.validateBrand(e); err != nil {
		return -1, err
	}

	if e.Id == 0 {
		return -1, errx.ErrInvalidInput
	}

	return s.repo.Update(ctx, e)
}

// Delete ...
func (s *brandService) Delete(ctx context.Context, id int32) (int64, error) {
	if id == 0 {
		return -1, errx.ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

// Count ...
func (s *brandService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// BulkInsertAll ...
func (s *brandService) BulkInsertAll(ctx context.Context, list []*domain.Brand) (int64, error) {
	if len(list) == 0 {
		return -1, errx.ErrInvalidInput
	}

	for _, brand := range list {
		if err := s.validateBrand(brand); err != nil {
			return -1, err
		}
	}

	return s.repo.BulkInsertAll(ctx, list)
}

// BulkInsert ...
func (s *brandService) BulkInsert(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error) {
	if len(list) == 0 {
		return -1, errx.ErrInvalidInput
	}

	for _, brand := range list {
		if err := s.validateBrand(brand); err != nil {
			return -1, err
		}
	}

	return s.repo.BulkInsert(ctx, list, batchSize)
}

// BulkUpdate ...
func (s *brandService) BulkUpdate(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error) {
	if len(list) == 0 {
		return -1, errx.ErrInvalidInput
	}

	for _, brand := range list {
		if err := s.validateBrand(brand); err != nil {
			return -1, err
		}
	}

	return s.repo.BulkUpdate(ctx, list, batchSize)
}
