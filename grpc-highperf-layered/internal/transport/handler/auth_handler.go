package handler

import (
	authv1 "codegen/api/gen/auth/v1"
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authHandler struct {
	authv1.UnimplementedAuthServiceServer
	logger *zerolog.Logger
}

func NewAuthHandler(logger *zerolog.Logger) *authHandler {
	server := &authHandler{}
	server.logger = logger

	return server
}

func (h *authHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AccessTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Register not implemented")
}

func (h *authHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AccessTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Login not implemented")
}

func (h *authHandler) ConfirmEmail(ctx context.Context, req *authv1.ConfirmEmailRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method ConfirmEmail not implemented")
}

func (h *authHandler) ForgotPassword(ctx context.Context, req *authv1.ForgotPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method ForgotPassword not implemented")
}

func (h *authHandler) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method ResetPassword not implemented")
}
