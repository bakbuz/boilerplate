package handler

import (
	catalogv1 "codegen/api/gen/catalog/v1"
	"codegen/internal/entity"
	"codegen/internal/service"
	"codegen/pkg/errx"
	"codegen/pkg/text"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type brandHandler struct {
	catalogv1.UnimplementedBrandServiceServer
	svc service.BrandService
}

func NewBrandHandler(svc service.BrandService) *brandHandler {
	return &brandHandler{svc: svc}
}

// ============================================================================
// MAPPER FUNCTIONS
// ============================================================================

// brandEntityToProto converts a Brand entity to protobuf Brand message
func brandEntityToProto(b *entity.Brand) *catalogv1.Brand {
	if b == nil {
		return nil
	}
	return &catalogv1.Brand{
		Id:        b.Id,
		Name:      b.Name,
		Slug:      b.Slug,
		Logo:      b.Logo,
		CreatedAt: timestamppb.New(b.CreatedAt),
	}
}

// brandCreateProtoToEntity converts CreateBrandRequest to Brand entity
func brandCreateProtoToEntity(req *catalogv1.CreateBrandRequest, currentUserId uuid.UUID) *entity.Brand {
	return &entity.Brand{
		Name:      req.Name,
		Slug:      text.Slugify(req.Name),
		Logo:      req.Logo,
		CreatedBy: currentUserId,
		CreatedAt: time.Now().UTC(),
	}
}

// brandUpdateProtoToEntity converts UpdateBrandRequest to Brand entity
func brandUpdateProtoToEntity(req *catalogv1.UpdateBrandRequest, currentUserId uuid.UUID) *entity.Brand {
	now := time.Now().UTC()
	return &entity.Brand{
		Id:        req.Id,
		Name:      req.Name,
		Slug:      text.Slugify(req.Name),
		Logo:      req.Logo,
		UpdatedBy: &currentUserId,
		UpdatedAt: &now,
	}
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// validateCreateBrandRequest validates a CreateBrandRequest
func (h *brandHandler) validateCreateBrandRequest(req *catalogv1.CreateBrandRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) > 255 {
		return status.Error(codes.InvalidArgument, "name too long (max 255 characters)")
	}
	return nil
}

// validateUpdateBrandRequest validates an UpdateBrandRequest
func (h *brandHandler) validateUpdateBrandRequest(req *catalogv1.UpdateBrandRequest) error {
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

func (h *brandHandler) List(ctx context.Context, req *catalogv1.ListBrandsRequest) (*catalogv1.ListBrandsResponse, error) {
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
		return &catalogv1.ListBrandsResponse{
			Total: 0,
			Items: []*catalogv1.Brand{},
		}, nil
	}

	// MAPPING: Entity List -> Proto List
	total := int32(len(list))
	protoItems := make([]*catalogv1.Brand, total)
	for i, b := range list {
		protoItems[i] = brandEntityToProto(b)
	}

	return &catalogv1.ListBrandsResponse{
		Total: total,
		Items: protoItems,
	}, nil
}

func (h *brandHandler) Get(ctx context.Context, req *catalogv1.GetBrandRequest) (*catalogv1.Brand, error) {
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

	return brandEntityToProto(item), nil
}

func (h *brandHandler) Create(ctx context.Context, req *catalogv1.CreateBrandRequest) (*catalogv1.Brand, error) {
	// Get current user Id from context
	currentUserId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user Id: %v", err)
	}

	// 1. Request Validation
	if err := h.validateCreateBrandRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity := brandCreateProtoToEntity(req, currentUserId)

	// 3. Service Call
	if err = h.svc.Create(ctx, domainEntity); err != nil {
		if errors.Is(err, errx.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid brand: %v", err)
		}
		if errors.Is(err, errx.ErrConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "brand already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create brand: %v", err)
	}

	// 4. MAPPING: Entity -> Proto Response
	return brandEntityToProto(domainEntity), nil
}

func (h *brandHandler) Update(ctx context.Context, req *catalogv1.UpdateBrandRequest) (*catalogv1.Brand, error) {
	// Get current user Id from context
	currentUserId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user Id: %v", err)
	}

	// 1. Request Validation
	if err := h.validateUpdateBrandRequest(req); err != nil {
		return nil, err
	}

	// 2. MAPPING: Proto -> Entity
	domainEntity := brandUpdateProtoToEntity(req, currentUserId)

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

	return brandEntityToProto(domainEntity), nil
}

func (h *brandHandler) Delete(ctx context.Context, req *catalogv1.DeleteBrandRequest) (*emptypb.Empty, error) {
	// Get current user Id from context
	_, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user Id: %v", err)
	}

	// Validate Id
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
