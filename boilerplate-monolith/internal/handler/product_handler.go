package handler

import (
	"codegen/api/pb"
	"codegen/internal/entity"
	"codegen/internal/service"
	"codegen/pkg/errx"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type productHandler struct {
	pb.UnimplementedCatalogServiceServer
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *productHandler {
	return &productHandler{svc: svc}
}

// ============================================================================
// MAPPER FUNCTIONS
// ============================================================================

// productEntityToProto converts a Product entity to protobuf Product message
func productEntityToProto(p *entity.Product) *pb.Product {
	if p == nil {
		return nil
	}

	product := &pb.Product{
		Id:        p.Id.String(),
		BrandId:   int32(p.BrandId),
		Name:      p.Name,
		Sku:       p.Sku,
		Summary:   p.Summary,
		Storyline: p.Storyline,
		Stock:     int32(p.StockQuantity),
		Price:     float32(p.Price),
		CreatedAt: timestamppb.New(p.CreatedAt),
	}

	return product
}

// productCreateProtoToEntity converts CreateProductRequest to Product entity
func productCreateProtoToEntity(req *pb.CreateProductRequest) (*entity.Product, error) {
	if req == nil {
		return nil, errx.ErrInvalidInput
	}

	newId, _ := uuid.NewV7()
	product := &entity.Product{
		Id:            newId,
		BrandId:       int(req.BrandId),
		Name:          req.Name,
		Sku:           req.Sku,
		Summary:       req.Summary,
		Storyline:     req.Storyline,
		StockQuantity: int(req.Stock),
		Price:         float64(req.Price),
	}

	return product, nil
}

// productUpdateProtoToEntity converts UpdateProductRequest to Product entity
func productUpdateProtoToEntity(req *pb.UpdateProductRequest) (*entity.Product, error) {
	if req == nil {
		return nil, errx.ErrInvalidInput
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid product id format")
	}

	product := &entity.Product{
		Id:            id,
		BrandId:       int(req.BrandId),
		Name:          req.Name,
		Sku:           req.Sku,
		Summary:       req.Summary,
		Storyline:     req.Storyline,
		StockQuantity: int(req.Stock),
		Price:         float64(req.Price),
	}

	return product, nil
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// validateCreateProductRequest validates a CreateProductRequest
func (h *productHandler) validateCreateProductRequest(req *pb.CreateProductRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "product is required")
	}

	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "product name is required")
	}
	if len(req.Name) > 255 {
		return status.Error(codes.InvalidArgument, "product name too long (max 255 characters)")
	}
	if req.BrandId == 0 {
		return status.Error(codes.InvalidArgument, "brand_id is required")
	}
	if req.Price < 0 {
		return status.Error(codes.InvalidArgument, "price must be non-negative")
	}
	if req.Stock < 0 {
		return status.Error(codes.InvalidArgument, "stock must be non-negative")
	}
	return nil
}

// validateUpdateProductRequest validates an UpdateProductRequest
func (h *productHandler) validateUpdateProductRequest(req *pb.UpdateProductRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "product is required")
	}
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "product id is required")
	}
	if _, err := uuid.Parse(req.Id); err != nil {
		return status.Error(codes.InvalidArgument, "invalid product id format")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "product name is required")
	}
	if len(req.Name) > 255 {
		return status.Error(codes.InvalidArgument, "product name too long (max 255 characters)")
	}
	if req.BrandId == 0 {
		return status.Error(codes.InvalidArgument, "brand_id is required")
	}
	if req.Price < 0 {
		return status.Error(codes.InvalidArgument, "price must be non-negative")
	}
	if req.Stock < 0 {
		return status.Error(codes.InvalidArgument, "stock must be non-negative")
	}
	return nil
}

// ============================================================================
// HANDLER METHODS
// ============================================================================

func (h *productHandler) GetProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// Context cancellation check
	if ctx.Err() != nil {
		return nil, status.Error(codes.Canceled, ctx.Err().Error())
	}

	list, err := h.svc.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}

	// Empty list optimization
	if len(list) == 0 {
		return &pb.ListProductsResponse{
			Total: 0,
			Items: []*pb.Product{},
		}, nil
	}

	// MAPPING: Entity List -> Proto List
	total := int32(len(list))
	protoItems := make([]*pb.Product, total)
	for i, p := range list {
		protoItems[i] = productEntityToProto(p)
	}

	return &pb.ListProductsResponse{
		Total: total,
		Items: protoItems,
	}, nil
}

func (h *productHandler) GetProduct(ctx context.Context, req *pb.ProductIdentifier) (*pb.GetProductResponse, error) {
	// Validate ID format
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id format")
	}

	item, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found: %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch product: %v", err)
	}

	return &pb.GetProductResponse{
		Product: productEntityToProto(item),
	}, nil
}

func (h *productHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductIdentifier, error) {
	// 1. Request Validation
	if err := h.validateCreateProductRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity, err := productCreateProtoToEntity(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product data: %v", err)
	}

	// 3. Service Call
	_, err = h.svc.Create(ctx, domainEntity)
	if err != nil {
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product: %v", err)
		}
		if errors.Is(err, errx.ErrConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "product already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	// 4. MAPPING: Entity -> Proto Response
	return &pb.ProductIdentifier{
		Id: domainEntity.Id.String(),
	}, nil
}

func (h *productHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*emptypb.Empty, error) {
	// 1. Request Validation
	if err := h.validateUpdateProductRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity, err := productUpdateProtoToEntity(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product data: %v", err)
	}

	// 3. Service Call
	rowsAffected, err := h.svc.Update(ctx, domainEntity)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found: %s", req.Id)
		}
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "product not found: %s", req.Id)
	}

	// 4. MAPPING: Entity -> Proto Response
	return &emptypb.Empty{}, nil
}

func (h *productHandler) DeleteProduct(ctx context.Context, req *pb.ProductIdentifier) (*emptypb.Empty, error) {
	// Validate ID format
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id format")
	}

	rowsAffected, err := h.svc.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found: %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "product not found: %s", req.Id)
	}

	return &emptypb.Empty{}, nil
}
