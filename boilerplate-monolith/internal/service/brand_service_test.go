package service

import (
	"codegen/internal/entity"
	"codegen/pkg/errx"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBrandRepository is a manual mock since we don't have mockgen output handy
type MockBrandRepository struct {
	mock.Mock
}

func (m *MockBrandRepository) GetAll(ctx context.Context) ([]*entity.Brand, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Brand), args.Error(1)
}

func (m *MockBrandRepository) GetByIds(ctx context.Context, ids []int32) ([]*entity.Brand, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Brand), args.Error(1)
}

func (m *MockBrandRepository) GetById(ctx context.Context, id int32) (*entity.Brand, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Brand), args.Error(1)
}

func (m *MockBrandRepository) Insert(ctx context.Context, e *entity.Brand) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockBrandRepository) Update(ctx context.Context, e *entity.Brand) (int64, error) {
	args := m.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) Delete(ctx context.Context, id int32) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) Upsert(ctx context.Context, e *entity.Brand) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockBrandRepository) BulkInsert(ctx context.Context, list []*entity.Brand) (int64, error) {
	args := m.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) BulkUpdate(ctx context.Context, list []*entity.Brand) (int64, error) {
	args := m.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBrandRepository) BulkInsertTran(ctx context.Context, list []*entity.Brand) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockBrandRepository) BulkUpdateTran(ctx context.Context, list []*entity.Brand) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func TestBrandService_GetAll(t *testing.T) {
	repo := new(MockBrandRepository)
	svc := NewBrandService(repo)
	ctx := context.Background()

	expectedApiBrands := []*entity.Brand{{Name: "b1"}, {Name: "b2"}}
	repo.On("GetAll", ctx).Return(expectedApiBrands, nil)

	list, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedApiBrands, list)
	repo.AssertExpectations(t)
}

func TestBrandService_GetById(t *testing.T) {
	repo := new(MockBrandRepository)
	svc := NewBrandService(repo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedBrand := &entity.Brand{Id: 1, Name: "b1"}
		repo.On("GetById", ctx, int32(1)).Return(expectedBrand, nil)

		brand, err := svc.GetById(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBrand, brand)
	})

	t.Run("invalid id", func(t *testing.T) {
		brand, err := svc.GetById(ctx, 0)
		assert.ErrorIs(t, err, errx.ErrInvalidInput)
		assert.Nil(t, brand)
	})

	t.Run("not found", func(t *testing.T) {
		repo.On("GetById", ctx, int32(2)).Return(nil, nil)
		brand, err := svc.GetById(ctx, 2)
		assert.ErrorIs(t, err, errx.ErrNotFound)
		assert.Nil(t, brand)
	})
}

func TestBrandService_Create(t *testing.T) {
	repo := new(MockBrandRepository)
	svc := NewBrandService(repo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		validBrand := &entity.Brand{Name: "valid", Slug: "valid"}
		repo.On("Insert", ctx, mock.MatchedBy(func(b *entity.Brand) bool {
			return b.Name == "valid" && !b.CreatedAt.IsZero()
		})).Return(nil)

		err := svc.Create(ctx, validBrand)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		invalidBrand := &entity.Brand{Name: "", Slug: "valid"}
		err := svc.Create(ctx, invalidBrand)
		assert.Error(t, err)
		assert.Equal(t, "brand name is required", err.Error())
	})
}

func TestBrandService_Update(t *testing.T) {
	repo := new(MockBrandRepository)
	svc := NewBrandService(repo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		validBrand := &entity.Brand{Id: 1, Name: "valid", Slug: "valid"}
		repo.On("Update", ctx, mock.MatchedBy(func(b *entity.Brand) bool {
			return b.UpdatedAt != nil && !b.UpdatedAt.IsZero()
		})).Return(int64(1), nil)

		rows, err := svc.Update(ctx, validBrand)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
	})

	t.Run("invalid input id 0", func(t *testing.T) {
		brand := &entity.Brand{Id: 0, Name: "valid", Slug: "valid"}
		rows, err := svc.Update(ctx, brand)
		assert.ErrorIs(t, err, errx.ErrInvalidInput)
		assert.Equal(t, int64(-1), rows)
	})
}

func TestBrandService_Validation(t *testing.T) {
	repo := new(MockBrandRepository)
	svc := NewBrandService(repo)
	ctx := context.Background()

	tests := []struct {
		name   string
		brand  *entity.Brand
		errMsg string
	}{
		{"empty name", &entity.Brand{Name: "", Slug: "slug"}, "brand name is required"},
		{"long name", &entity.Brand{Name: strings.Repeat("a", 256), Slug: "slug"}, "brand name must not exceed 255 characters"},
		{"long name multibyte", &entity.Brand{Name: strings.Repeat("日", 256), Slug: "slug"}, "brand name must not exceed 255 characters"},
		{"name multibyte ok", &entity.Brand{Name: strings.Repeat("日", 255), Slug: "slug"}, ""},
		{"empty slug", &entity.Brand{Name: "name", Slug: ""}, "brand slug is required"},
		{"long slug", &entity.Brand{Name: "name", Slug: strings.Repeat("a", 256)}, "brand slug must not exceed 255 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errMsg == "" {
				repo.On("Insert", ctx, mock.Anything).Return(nil)
			}

			err := svc.Create(ctx, tt.brand)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
