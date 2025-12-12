package handler

import (
	"codegen/api/pb"
	"codegen/internal/entity"
	"codegen/internal/service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler interface {
	GetProducts(context.Context, *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	GetProduct(context.Context, *pb.ProductIdentifier) (*pb.GetProductResponse, error)
	CreateProduct(context.Context, *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	UpdateProduct(context.Context, *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
	DeleteProduct(context.Context, *pb.ProductIdentifier) (*pb.SuccessResponse, error)

	ListBrands(context.Context, *pb.ListBrandsRequest) (*pb.ListBrandsResponse, error)
	GetBrand(context.Context, *pb.BrandIdentifier) (*pb.GetBrandResponse, error)
	CreateBrand(context.Context, *pb.CreateBrandRequest) (*pb.BrandIdentifier, error)
	UpdateBrand(context.Context, *pb.UpdateBrandRequest) (*pb.SuccessResponse, error)
	DeleteBrand(context.Context, *pb.BrandIdentifier) (*pb.SuccessResponse, error)
}

type productHandler struct {
	pb.UnimplementedCatalogServiceServer
	psvc service.ProductService
	bsvc service.BrandService
}

func NewProductHandler(psvc service.ProductService, bsvc service.BrandService) ProductHandler {
	return &productHandler{psvc: psvc, bsvc: bsvc}
}

/*
func (h *productHandler) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.CreateProductResponse, error) {
	p := &service.EntityProductFromProto(req.Product) // convert without allocations if possible
	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "create failed: %v", err)
	}
	return &gen.CreateProductResponse{Product: service.ProtoFromEntity(p)}, nil
}*/

func (h *productHandler) GetProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetProducts not implemented")
}
func (h *productHandler) GetProduct(ctx context.Context, req *pb.ProductIdentifier) (*pb.GetProductResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetProduct not implemented")
}
func (h *productHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateProduct not implemented")
}
func (h *productHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateProduct not implemented")
}
func (h *productHandler) DeleteProduct(ctx context.Context, req *pb.ProductIdentifier) (*pb.SuccessResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteProduct not implemented")
}

func (h *productHandler) ListBrands(ctx context.Context, req *pb.ListBrandsRequest) (*pb.ListBrandsResponse, error) {
	list, err := h.bsvc.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch brands: %v", err)
	}

	// MAPPING: Entity List -> Proto List
	var total int32 = int32(len(list))
	var protoItems []*pb.Brand = make([]*pb.Brand, total)
	for i, b := range list {
		protoItems[i] = &pb.Brand{
			Id:   b.Id,
			Name: b.Name,
			Slug: b.Slug,
			Logo: b.Logo,
		}
	}

	return &pb.ListBrandsResponse{
		Total: total,
		Items: protoItems,
	}, nil
}

func (h *productHandler) GetBrand(ctx context.Context, req *pb.BrandIdentifier) (*pb.GetBrandResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetBrand not implemented")
}

func (h *productHandler) CreateBrand(ctx context.Context, req *pb.CreateBrandRequest) (*pb.BrandIdentifier, error) {
	// 1. Request Validation (Stateless - Yapısal Doğrulama)
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// 2. MAPPING: Proto -> Entity (Bunu ayrı bir mapper fonksiyonuna da taşıyabilirsin)
	domainEntity := &entity.Brand{
		Name: req.Name,
		Slug: req.Slug,
		Logo: req.Logo,
	}

	// 3. Service Call (Saf Go struct ile)
	err := h.bsvc.Create(ctx, domainEntity)
	if err != nil {
		// Hata yönetimi ve Loglama burada yapılır
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	// 4. MAPPING: Entity -> Proto Response
	return &pb.BrandIdentifier{
		Id: domainEntity.Id,
	}, nil
}

func (h *productHandler) UpdateBrand(ctx context.Context, req *pb.UpdateBrandRequest) (*pb.SuccessResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateBrand not implemented")
}

func (h *productHandler) DeleteBrand(ctx context.Context, req *pb.BrandIdentifier) (*pb.SuccessResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteBrand not implemented")
}
