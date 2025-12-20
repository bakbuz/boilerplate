package benchmark_test

import (
	"codegen/internal/domain"
	"codegen/internal/service"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type brandRepositoryMock struct {
	mock.Mock
}

// Ping implements database.Repository
func (r *brandRepositoryMock) Ping(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}

func (r *brandRepositoryMock) GetAll(ctx context.Context) ([]*domain.Brand, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*domain.Brand), args.Error(1)
}

func (r *brandRepositoryMock) GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).([]*domain.Brand), args.Error(1)
}

func (r *brandRepositoryMock) GetById(ctx context.Context, id int32) (*domain.Brand, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*domain.Brand), args.Error(1)
}

func (r *brandRepositoryMock) Insert(ctx context.Context, e *domain.Brand) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *brandRepositoryMock) Update(ctx context.Context, e *domain.Brand) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) Delete(ctx context.Context, id int32) (int64, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) SoftDelete(ctx context.Context, id int32, deletedBy uuid.UUID) (int64, error) {
	args := r.Called(ctx, id, deletedBy)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) Count(ctx context.Context) (int64, error) {
	args := r.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) Upsert(ctx context.Context, e *domain.Brand) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *brandRepositoryMock) BulkInsertOneShot(ctx context.Context, list []*domain.Brand) (int64, error) {
	args := r.Called(ctx, list)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) BulkInsert(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error) {
	args := r.Called(ctx, list, batchSize)
	return args.Get(0).(int64), args.Error(1)
}

func (r *brandRepositoryMock) BulkUpdate(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error) {
	args := r.Called(ctx, list, batchSize)
	return args.Get(0).(int64), args.Error(1)
}

func BenchmarkBrandGetLists(b *testing.B) {
	r := new(brandRepositoryMock)
	s := service.NewBrandService(r)

	r.On("GetAll", mock.Anything).Return([]*domain.Brand{
		{
			Id:        1,
			Name:      "list1",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        2,
			Name:      "list2",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        3,
			Name:      "list3",
			CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
		},
		{
			Id:        4,
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
