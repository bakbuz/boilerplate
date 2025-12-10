package grpctransport

import (
	"codegen/api/gen/codegen/api/gen"
	"codegen/internal/service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productHandler struct {
	svc *service.ProductService
	// Use sync.Pool for reusable buffers if needed for marshalling
}

func NewProductHandler(svc *service.ProductService) gen.ProductServiceServer {
	return &productHandler{svc: svc}
}

func (h *productHandler) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.CreateProductResponse, error) {
	p := &service.EntityProductFromProto(req.Product) // convert without allocations if possible
	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "create failed: %v", err)
	}
	return &gen.CreateProductResponse{Product: service.ProtoFromEntity(p)}, nil
}
