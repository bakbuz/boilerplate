package service_test

import (
	"codegen/internal/domain"
	"codegen/internal/service"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func (r *repositoryMock) GetAll(ctx context.Context) ([]*domain.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (r *repositoryMock) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (r *repositoryMock) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (r *repositoryMock) Insert(ctx context.Context, e *domain.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Update(ctx context.Context, e *domain.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	args := r.Called(ctx, id, deletedBy)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Count(ctx context.Context) (int64, error) {
	args := r.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Upsert(ctx context.Context, e *domain.Product) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *repositoryMock) BulkInsert(ctx context.Context, list []*domain.Product) (int64, error) {
	args := r.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) BulkUpdate(ctx context.Context, list []*domain.Product) (int64, error) {
	args := r.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(*domain.ProductSearchResult), args.Error(1)
}

func BenchmarkGetLists(b *testing.B) {
	r := new(repositoryMock)
	s := service.NewProductService(r)

	r.On("GetAll", mock.Anything).Return([]*domain.Product{
		{
			Id:        uuid.New(),
			Name:      "list1",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        uuid.New(),
			Name:      "list2",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        uuid.New(),
			Name:      "list3",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        uuid.New(),
			Name:      "list4",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetAll(context.Background())
	}
}
