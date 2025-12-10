package service_test

import (
	"context"
	"testing"
	"time"

	"codegen/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func (r *repositoryMock) GetAll(ctx context.Context) ([]*entity.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (r *repositoryMock) GetByIds(ctx context.Context) ([]*entity.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (r *repositoryMock) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (r *repositoryMock) Insert(ctx context.Context, e *entity.Product) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *repositoryMock) Update(ctx context.Context, e *entity.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy any) (int64, error) {
	args := r.Called(ctx, id, deletedBy)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) Count(ctx context.Context) (int64, error) {
	args := r.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

//var _ repository.TodoRepository = (repositoryMock)(nil)

func BenchmarkGetLists(b *testing.B) {
	r := new(repositoryMock)
	//s := service.NewTodoService(&r)

	r.On("GetLists", mock.Anything, mock.Anything).Return([]*entity.Product{
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
		r.GetAll(context.Background())
	}
}
