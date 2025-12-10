package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/yourusername/grpc-highperf-backend/internal/domain"
    "github.com/yourusername/grpc-highperf-backend/internal/service"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, product *domain.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}

func TestProductService_Create(t *testing.T) {
    mockRepo := new(MockRepository)
    svc := service.NewProductService(mockRepo)
    
    product := &domain.Product{
        Name:  "Test Product",
        Price: 99.99,
        Stock: 10,
    }
    
    mockRepo.On("Create", mock.Anything, product).Return(nil)
    
    ctx := context.Background()
    err := svc.Create(ctx, product)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}