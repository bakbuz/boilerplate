package interceptor

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	BearerPrefix string = "Bearer "
	userIdKey    string = "user_id"
)

var publicEndpoints []string = []string{
	"/api.v1.AuthService/Login",
	"/api.v1.AuthService/Register",
	"/grpc.health.v1.Health/Check",
}

func AuthInterceptor(jwtSecretKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Public endpoints
		if isPublicEndpoint(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		tokenString := strings.TrimPrefix(tokens[0], BearerPrefix)
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}

		// Validate token
		claims, err := validateToken(jwtSecretKey, tokenString)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add claims to context
		userId := claims["sub"].(string)
		ctx = context.WithValue(ctx, userIdKey, userId)

		return handler(ctx, req)
	}
}

func isPublicEndpoint(method string) bool {
	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}
	return false
}

func validateToken(jwtSecretKey, tokenString string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}
	claims := t.Claims.(jwt.MapClaims)
	//sub := claims["sub"].(string)
	return claims, nil
}
