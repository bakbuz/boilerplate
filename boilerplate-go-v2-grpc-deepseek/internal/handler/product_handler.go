package handler

import (
    "context"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    apiv1 "github.com/yourusername/grpc-highperf-backend/api/v1"
    "github.com/yourusername/grpc-highperf-backend/internal/domain"
    "github.com/yourusername/grpc-highperf-backend/internal/service"
    "github.com/yourusername/grpc-highperf-backend/internal/pkg/errors"
)

type productHandler struct {
    apiv1.UnimplementedProductServiceServer
    service service.ProductService
}

// Zero-allocation handler using sync.Pool for request/response objects
var productPool = &sync.Pool{
    New: func() interface{} {
        return &apiv1.ProductResponse{}
    },
}

func RegisterProductHandler(s *grpc.Server, svc service.ProductService) {
    handler := &productHandler{service: svc}
    apiv1.RegisterProductServiceServer(s, handler)
}

func (h *productHandler) CreateProduct(ctx context.Context, req *apiv1.CreateProductRequest) (*apiv1.ProductResponse, error) {
    // Input validation
    if req.Name == "" || req.Price <= 0 {
        return nil, status.Error(codes.InvalidArgument, "invalid product data")
    }

    product := &domain.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
    }

    created, err := h.service.Create(ctx, product)
    if err != nil {
        return nil, errors.ToGRPCError(err)
    }

    // Reuse response object from pool
    resp := productPool.Get().(*apiv1.ProductResponse)
    defer productPool.Put(resp)

    // Reset and populate response
    resp.Reset()
    resp.Product = &apiv1.Product{
        Id:          created.ID,
        Name:        created.Name,
        Description: created.Description,
        Price:       created.Price,
        Stock:       created.Stock,
        CreatedAt:   timestamppb.New(created.CreatedAt),
        UpdatedAt:   timestamppb.New(created.UpdatedAt),
    }

    return resp, nil
}

func (h *productHandler) GetProduct(ctx context.Context, req *apiv1.GetProductRequest) (*apiv1.ProductResponse, error) {
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "product id is required")
    }

    product, err := h.service.GetByID(ctx, req.Id)
    if err != nil {
        return nil, errors.ToGRPCError(err)
    }

    resp := productPool.Get().(*apiv1.ProductResponse)
    defer productPool.Put(resp)

    resp.Reset()
    resp.Product = &apiv1.Product{
        Id:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
        CreatedAt:   timestamppb.New(product.CreatedAt),
        UpdatedAt:   timestamppb.New(product.UpdatedAt),
    }

    return resp, nil
}

// Optimized ListProducts with pagination
func (h *productHandler) ListProducts(ctx context.Context, req *apiv1.ListProductsRequest) (*apiv1.ListProductsResponse, error) {
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 || req.PageSize > 100 {
        req.PageSize = 20
    }

    products, total, err := h.service.List(ctx, req.Page, req.PageSize, req.Filter)
    if err != nil {
        return nil, errors.ToGRPCError(err)
    }

    resp := &apiv1.ListProductsResponse{
        Products: make([]*apiv1.Product, 0, len(products)),
        Total:    int32(total),
        Page:     int32(req.Page),
        PageSize: int32(req.PageSize),
    }

    // Pre-allocate and populate response
    for _, p := range products {
        resp.Products = append(resp.Products, &apiv1.Product{
            Id:          p.ID,
            Name:        p.Name,
            Description: p.Description,
            Price:       p.Price,
            Stock:       p.Stock,
            CreatedAt:   timestamppb.New(p.CreatedAt),
            UpdatedAt:   timestamppb.New(p.UpdatedAt),
        })
    }

    return resp, nil
}