package service

import (
	"codegen/internal/entity"
	"codegen/internal/repository/dto"
	"codegen/pkg/errx"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a manual mock
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) Insert(ctx context.Context, e *entity.Product) (int64, error) {
	args := m.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, e *entity.Product) (int64, error) {
	args := m.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	args := m.Called(ctx, id, deletedBy)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) Upsert(ctx context.Context, e *entity.Product) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockProductRepository) BulkInsert(ctx context.Context, list []*entity.Product) (int64, error) {
	args := m.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) BulkUpdate(ctx context.Context, list []*entity.Product) (int64, error) {
	args := m.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) Search(ctx context.Context, filter *dto.ProductSearchFilter) (*dto.ProductSearchResult, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductSearchResult), args.Error(1)
}

func TestProductService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		validProduct := &entity.Product{
			Name:          "Valid Product",
			BrandId:       1,
			StockQuantity: 10,
			Price:         99.99,
		}

		repo.On("Insert", ctx, mock.MatchedBy(func(p *entity.Product) bool {
			return p.Name == "Valid Product" && p.BrandId == 1 && p.Id != uuid.Nil
		})).Return(int64(1), nil)

		rows, err := svc.Create(ctx, validProduct)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
		repo.AssertExpectations(t)
	})

	t.Run("validation error - empty name", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		invalidProduct := &entity.Product{
			Name:    "   ", // Should be trimmed to empty
			BrandId: 1,
		}

		rows, err := svc.Create(ctx, invalidProduct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product name is required")
		assert.Equal(t, int64(-1), rows)
	})

	t.Run("validation error - negative price", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		invalidProduct := &entity.Product{
			Name:    "Valid Name",
			BrandId: 1,
			Price:   -10,
		}

		rows, err := svc.Create(ctx, invalidProduct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "price cannot be negative")
		assert.Equal(t, int64(-1), rows)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		validProduct := &entity.Product{
			Name:    "Repo Error",
			BrandId: 1,
		}

		repo.On("Insert", ctx, mock.Anything).Return(int64(0), errx.ErrConflict)

		rows, err := svc.Create(ctx, validProduct)
		assert.ErrorIs(t, err, errx.ErrConflict)
		assert.Equal(t, int64(0), rows)
	})
}

func TestProductService_Update(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		validProduct := &entity.Product{
			Id:      id,
			Name:    "Updated Product",
			BrandId: 1,
		}

		repo.On("Update", ctx, mock.MatchedBy(func(p *entity.Product) bool {
			return p.Name == "Updated Product" && p.Id == id
		})).Return(int64(1), nil)

		rows, err := svc.Update(ctx, validProduct)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
	})

	t.Run("validation error - invalid id", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		invalidProduct := &entity.Product{
			Id:      uuid.Nil,
			Name:    "Valid Name",
			BrandId: 1,
		}

		rows, err := svc.Update(ctx, invalidProduct)
		assert.ErrorIs(t, err, errx.ErrInvalidInput)
		assert.Equal(t, int64(-1), rows)
	})
}

func TestProductService_GetById(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		expectedProduct := &entity.Product{Id: id, Name: "Found"}
		repo.On("GetById", ctx, id).Return(expectedProduct, nil)

		product, err := svc.GetById(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		repo.On("GetById", ctx, id).Return(nil, nil)

		product, err := svc.GetById(ctx, id)
		assert.ErrorIs(t, err, errx.ErrNotFound)
		assert.Nil(t, product)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		product, err := svc.GetById(ctx, uuid.Nil)
		assert.ErrorIs(t, err, errx.ErrInvalidInput)
		assert.Nil(t, product)
	})
}

func TestProductService_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		filter := &dto.ProductSearchFilter{Take: 10, Skip: 0}
		expectedResult := &dto.ProductSearchResult{Total: 1, Items: []*entity.Product{{Name: "Search Result"}}}
		repo.On("Search", ctx, mock.MatchedBy(func(f *dto.ProductSearchFilter) bool {
			return f.Take == 10
		})).Return(expectedResult, nil)

		result, err := svc.Search(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("default pagination", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		filter := &dto.ProductSearchFilter{} // Take 0
		repo.On("Search", ctx, mock.MatchedBy(func(f *dto.ProductSearchFilter) bool {
			return f.Take == 50 // Should default to 50
		})).Return(&dto.ProductSearchResult{}, nil)

		_, err := svc.Search(ctx, filter)
		assert.NoError(t, err)
	})

	t.Run("invalid pagination", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		filter := &dto.ProductSearchFilter{Take: 1001}
		result, err := svc.Search(ctx, filter)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "take parameter must not exceed 1000")
		assert.Nil(t, result)
	})
}

func TestProductService_Sanitization(t *testing.T) {
	ctx := context.Background()

	t.Run("trim strings", func(t *testing.T) {
		repo := new(MockProductRepository)
		svc := NewProductService(repo)
		sku := "  SKU-123  "
		input := &entity.Product{
			Name:    "  Trim Me  ",
			Sku:     &sku,
			BrandId: 1,
		}

		repo.On("Insert", ctx, mock.MatchedBy(func(p *entity.Product) bool {
			return p.Name == "Trim Me" && *p.Sku == "SKU-123"
		})).Return(int64(1), nil)

		_, err := svc.Create(ctx, input)
		assert.NoError(t, err)
	})
}
