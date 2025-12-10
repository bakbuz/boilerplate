package service

import (
	"codegen/internal/errorcodes"
	"codegen/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepo
}

func NewOrderService(r repository.ProductRepository) *OrderService {
	return &OrderService{repo: r}
}

func (s *OrderService) PlaceOrder(userID string, productID string, quantity int) error {
	// validation
	if quantity <= 0 {
		return errorcodes.ErrInvalidInput
	}
	return s.repo.Create(ctx, p)
}
