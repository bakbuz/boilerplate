package handler

import (
	"codegen/api/pb"
	"codegen/internal/entity"
	"codegen/internal/service"
	"codegen/pkg/errx"
	"codegen/pkg/text"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type brandHandler struct {
	pb.UnimplementedBrandServiceServer
	svc service.BrandService
}

func NewBrandHandler(svc service.BrandService) *brandHandler {
	return &brandHandler{svc: svc}
}

// ============================================================================
// MAPPER FUNCTIONS
// ============================================================================

// brandEntityToProto converts a Brand entity to protobuf Brand message
func brandEntityToProto(b *entity.Brand) *pb.Brand {
	if b == nil {
		return nil
	}
	return &pb.Brand{
		Id:        b.Id,
		Name:      b.Name,
		Slug:      b.Slug,
		Logo:      b.Logo,
		CreatedAt: timestamppb.New(b.CreatedAt),
	}
}

// brandCreateProtoToEntity converts CreateBrandRequest to Brand entity
func brandCreateProtoToEntity(req *pb.CreateBrandRequest) *entity.Brand {
	return &entity.Brand{
		Name: req.Name,
		Slug: text.Slugify(req.Name),
		Logo: req.Logo,
	}
}

// brandUpdateProtoToEntity converts UpdateBrandRequest to Brand entity
func brandUpdateProtoToEntity(req *pb.UpdateBrandRequest) *entity.Brand {
	return &entity.Brand{
		Id:   req.Id,
		Name: req.Name,
		Slug: text.Slugify(req.Name),
		Logo: req.Logo,
	}
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// validateCreateBrandRequest validates a CreateBrandRequest
func (h *brandHandler) validateCreateBrandRequest(req *pb.CreateBrandRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) > 255 {
		return status.Error(codes.InvalidArgument, "name too long (max 255 characters)")
	}
	return nil
}

// validateUpdateBrandRequest validates an UpdateBrandRequest
func (h *brandHandler) validateUpdateBrandRequest(req *pb.UpdateBrandRequest) error {
	if req.GetId() == 0 {
		return status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) > 255 {
		return status.Error(codes.InvalidArgument, "name too long (max 255 characters)")
	}
	return nil
}

// ============================================================================
// HANDLER METHODS
// ============================================================================

func (h *brandHandler) ListBrands(ctx context.Context, req *emptypb.Empty) (*pb.ListBrandsResponse, error) {
	// Context cancellation check
	if ctx.Err() != nil {
		return nil, status.Error(codes.Canceled, ctx.Err().Error())
	}

	list, err := h.svc.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch brands: %v", err)
	}

	// Empty list optimization
	if len(list) == 0 {
		return &pb.ListBrandsResponse{
			Total: 0,
			Items: []*pb.Brand{},
		}, nil
	}

	// MAPPING: Entity List -> Proto List
	total := int32(len(list))
	protoItems := make([]*pb.Brand, total)
	for i, b := range list {
		protoItems[i] = brandEntityToProto(b)
	}

	return &pb.ListBrandsResponse{
		Total: total,
		Items: protoItems,
	}, nil
}

func (h *brandHandler) GetBrand(ctx context.Context, req *pb.BrandIdentifier) (*pb.GetBrandResponse, error) {
	item, err := h.svc.GetById(ctx, req.Id)
	if err != nil {
		// Distinguish between NotFound and other errors
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "brand not found: %d", req.Id)
		}
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, "invalid brand id")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch brand: %v", err)
	}

	return &pb.GetBrandResponse{
		Brand: brandEntityToProto(item),
	}, nil
}

func (h *brandHandler) CreateBrand(ctx context.Context, req *pb.CreateBrandRequest) (*pb.BrandIdentifier, error) {
	// 1. Request Validation
	if err := h.validateCreateBrandRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity := brandCreateProtoToEntity(req)

	// 3. Service Call
	err := h.svc.Create(ctx, domainEntity)
	if err != nil {
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid brand: %v", err)
		}
		if errors.Is(err, errx.ErrConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "brand already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create brand: %v", err)
	}

	// 4. MAPPING: Entity -> Proto Response
	return &pb.BrandIdentifier{
		Id: domainEntity.Id,
	}, nil
}

func (h *brandHandler) UpdateBrand(ctx context.Context, req *pb.UpdateBrandRequest) (*emptypb.Empty, error) {
	// 1. Request Validation
	if err := h.validateUpdateBrandRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity := brandUpdateProtoToEntity(req)

	// 3. Service Call
	rowsAffected, err := h.svc.Update(ctx, domainEntity)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "brand not found: %d", req.Id)
		}
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid brand: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update brand: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brand not found: %d", req.Id)
	}

	return &emptypb.Empty{}, nil
}

func (h *brandHandler) DeleteBrand(ctx context.Context, req *pb.BrandIdentifier) (*emptypb.Empty, error) {
	// Validate ID
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "brand id is required")
	}

	// Service Call
	rowsAffected, err := h.svc.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "brand not found: %d", req.Id)
		}
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, "invalid brand id")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete brand: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brand not found: %d", req.Id)
	}

	return &emptypb.Empty{}, nil
}
