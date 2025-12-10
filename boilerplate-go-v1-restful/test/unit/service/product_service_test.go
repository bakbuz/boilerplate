package service_test

import (
	"context"
	"testing"

	"codegen/internal/entity"
	"codegen/internal/repository/dto"
	"codegen/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ProductRepository Mock ...
type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Search(ctx context.Context, filter *dto.ProductFilter) (int, []*entity.Product, error) {
	args := r.Called(ctx, filter)
	return args.Int(0), args.Get(0).([]*entity.Product), args.Error(1)
}
func (r *repositoryMock) GetAll(ctx context.Context) ([]*entity.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*entity.Product), args.Error(1)
}
func (r *repositoryMock) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).([]*entity.Product), args.Error(1)
}
func (r *repositoryMock) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*entity.Product), args.Error(1)
}
func (r *repositoryMock) Insert(ctx context.Context, e *entity.Product) (int64, error) {
	args := r.Called(ctx, e)
	return args.Get(0).(int64), args.Error(1)
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
func (r *repositoryMock) BulkInsert(ctx context.Context, list []*entity.Product) error {
	args := r.Called(ctx, list)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	r := &repositoryMock{}
	s := service.NewProductService(r)

	e := &entity.Product{}
	r.On("CreateProduct", mock.Anything, e).Return(1, nil)

	s.CreateProduct(context.Background(), e)
	r.AssertCalled(t, "CreateProduct", mock.Anything, e)
	r.AssertNumberOfCalls(t, "CreateProduct", 1)
}

/*
func TestAddTask(t *testing.T) {
	input := []struct {
		name          string
		dbResult      int
		dbError       error
		expected      int
		expectedError error
	}{
		{
			"success",
			1,
			nil,
			1,
			nil,
		},
		{
			"product not found error",
			-1,
			repository.ErrProductForeignKeyViolation,
			-1,
			service.ErrProductNotFound,
		},
		{
			"unknown error",
			-1,
			fmt.Errorf("unknown error"),
			-1,
			fmt.Errorf("unknown error"),
		},
	}

	for _, tc := range input {
		t.Run(tc.name, func(t *testing.T) {
			r := &repositoryMock{}
			s := service.NewProductService(r)

			r.On("CreateTask", mock.Anything, 1, "task1").Return(tc.dbResult, tc.dbError)

			res, err := s.AddTask(context.Background(), 1, "task1")

			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCloseTask(t *testing.T) {
	input := []struct {
		name          string
		dbError       error
		expectedError error
	}{
		{
			"success",
			nil,
			nil,
		},
		{
			"unknown error",
			fmt.Errorf("unknown error"),
			fmt.Errorf("unknown error"),
		},
	}

	for _, tc := range input {
		t.Run(tc.name, func(t *testing.T) {
			r := &repositoryMock{}
			s := service.NewProductService(r)

			r.On("DeleteTask", mock.Anything, 1).Return(tc.dbError)

			err := s.CloseTask(context.Background(), 1)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetProducts(t *testing.T) {
	testCases := []struct {
		name          string
		dbResult      []*entity.Product
		dbError       error
		expected      []*entity.Product
		expectedError error
	}{
		{
			"success",
			[]*entity.Product{
				{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				{
					Id:        2,
					Name:      "product2",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
			},
			nil,
			[]*entity.Product{
				{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				{
					Id:        2,
					Name:      "product2",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &repositoryMock{}
			s := service.NewProductService(r)

			r.On("GetProducts", mock.Anything).Return(tc.dbResult, tc.dbError)

			products, err := s.GetProducts(context.Background())
			require.NoError(t, err)

			assert.Equal(t, tc.expected, products)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetProductWithTasks(t *testing.T) {

	testCases := []struct {
		name          string
		dbResult      *dto.ProductWithTasksDto
		dbError       error
		expected      *dto.ProductWithTasksDto
		expectedError error
	}{
		{
			"success",
			&dto.ProductWithTasksDto{
				Product: &entity.Product{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				Tasks: []*entity.Task{
					{
						Id:        1,
						ProductId: 1,
						Name:      "task1",
						Done:      false,
						CreatedAt: time.Date(2021, 2, 3, 4, 5, 6, 7, time.UTC),
					},
					{
						Id:        2,
						ProductId: 1,
						Name:      "task2",
						Done:      false,
						CreatedAt: time.Date(2021, 3, 4, 5, 6, 7, 8, time.UTC),
					},
				},
			},
			nil, //dbError
			&dto.ProductWithTasksDto{
				Product: &entity.Product{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				Tasks: []*entity.Task{
					{
						Id:        1,
						ProductId: 1,
						Name:      "task1",
						Done:      false,
						CreatedAt: time.Date(2021, 2, 3, 4, 5, 6, 7, time.UTC),
					},
					{
						Id:        2,
						ProductId: 1,
						Name:      "task2",
						Done:      false,
						CreatedAt: time.Date(2021, 3, 4, 5, 6, 7, 8, time.UTC),
					},
				},
			},
			nil, // expectedError
		},
		{
			"success without tasks",
			&dto.ProductWithTasksDto{
				Product: &entity.Product{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				Tasks: nil,
			},
			nil,
			&dto.ProductWithTasksDto{
				Product: &entity.Product{
					Id:        1,
					Name:      "product1",
					CreatedAt: time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC),
				},
				Tasks: nil,
			},
			nil,
		},
		{
			"product not found error #1",
			&dto.ProductWithTasksDto{},
			nil,
			&dto.ProductWithTasksDto{},
			nil,
		},
		{
			"unknown error",
			nil,
			fmt.Errorf("unknown error"),
			nil,
			fmt.Errorf("unknown error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &repositoryMock{}
			s := service.NewProductService(r)

			r.On("GetProductWithTasks", mock.Anything, 1).Return(tc.dbResult, tc.dbError)

			res, err := s.GetProductWithTasks(context.Background(), 1)

			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
*/
