package grpctransport

import (
	"codegen/api/gen/codegen/api/gen"
	"codegen/internal/service"
	"context"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	grpc *grpc.Server
}

func NewServer(ps *service.ProductService, os *service.OrderService, opts ...grpc.ServerOption) *Server {
	gs := grpc.NewServer(opts...)
	gen.RegisterProductServiceServer(gs, NewProductHandler(ps))
	gen.RegisterOrderServiceServer(gs, NewOrderHandler(os))
	return &Server{grpc: gs}
}

func (s *Server) Serve(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	return s.grpc.Serve(lis)
}

func (s *Server) GracefulStop(ctx context.Context) {
	s.grpc.GracefulStop()
}
