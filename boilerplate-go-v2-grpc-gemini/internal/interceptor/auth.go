package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Public endpointlere izin verilebilir (örn: Login)
		// if info.FullMethod == "/catalog.AuthService/Login" { return handler(ctx, req) }

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := strings.TrimPrefix(values[0], "Bearer ")

		// Burada Token doğrulama işlemi yapılır (hmac/rsa)
		valid := validateToken(token, jwtSecret)
		if !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Token geçerliyse devam et
		return handler(ctx, req)
	}
}

func validateToken(token, jwtSecret string) bool {
	panic("unimplemented")
}
