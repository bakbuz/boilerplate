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

type productRepositoryMock struct {
	mock.Mock
}

// Ping implements database.Repository
func (r *productRepositoryMock) Ping(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}

func (r *productRepositoryMock) GetAll(ctx context.Context) ([]*domain.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (r *productRepositoryMock) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (r *productRepositoryMock) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (r *productRepositoryMock) Insert(ctx context.Context, e *domain.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) Update(ctx context.Context, e *domain.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	args := r.Called(ctx, id, deletedBy)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) Count(ctx context.Context) (int64, error) {
	args := r.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) Upsert(ctx context.Context, e *domain.Product) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *productRepositoryMock) BulkInsertOneShot(ctx context.Context, list []*domain.Product) (int64, error) {
	args := r.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) BulkInsert(ctx context.Context, list []*domain.Product, batchSize int) (int64, error) {
	args := r.Called(ctx, list, batchSize)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) BulkUpdate(ctx context.Context, list []*domain.Product, batchSize int) (int64, error) {
	args := r.Called(ctx, list, batchSize)
	return args.Get(0).(int64), args.Error(1)
}

func (r *productRepositoryMock) Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(*domain.ProductSearchResult), args.Error(1)
}

func BenchmarkProductGetLists(b *testing.B) {
	r := new(productRepositoryMock)
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
		_, err := s.GetAll(context.Background())
		if err != nil {
			b.Error(err)
		}
	}
}
