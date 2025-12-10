package unit_test

import (
	"codegen/internal/entity"
	"codegen/internal/service"
	"context"
	"testing"
)

func TestProductCreate_Invalid(t *testing.T) {
	repo := &mockProductRepo{}
	svc := service.NewProductService(repo)
	p := &entity.Product{Name: "", Price: 0}
	err := svc.Create(context.Background(), p)
	if err != service.ErrInvalidInput {
		t.Fatalf("expected invalid input")
	}
}
