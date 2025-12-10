package service

import (
	"codegen/internal/repository"
	"codegen/proto/pb"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Hataları google.golang.org/grpc/status ile yönetiyoruz.

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	repo repository.Repository
}

func NewProductServer(repo repository.Repository) *ProductServer {
	return &ProductServer{repo: repo}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	product, err := s.repo.GetProductByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		},
	}, nil
}
