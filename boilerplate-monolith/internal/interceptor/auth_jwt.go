package interceptor

import (
	"codegen/api/pb"
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Context key çakışmalarını önlemek için özel tip
type contextKey string

const (
	BearerPrefix string     = "Bearer "
	UserIDKey    contextKey = "user_id"
)

// Hızlı arama için map kullanıyoruz (Set mantığı)
// pb paketini burada kullanmak çok mantıklı değil ama olası endpoint değişikliklerini önlemek için kullanıyoruz
var publicEndpoints = map[string]struct{}{
	pb.DemoService_ListDemos_FullMethodName: {},
	pb.DemoService_GetDemo_FullMethodName:   {},
	pb.DemoService_Create_FullMethodName:    {},
	"/api.v1.AuthService/Login":             {},
	"/api.v1.AuthService/Register":          {},
	"/grpc.health.v1.Health/Check":          {},
}

func AuthInterceptor(jwtSecretKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Public endpoint kontrolü (O(1) complexity)
		if _, ok := publicEndpoints[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// 2. Metadata kontrolü
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		// 3. Token format kontrolü
		tokenString := strings.TrimPrefix(tokens[0], BearerPrefix)
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}

		// 4. Token doğrulama ve Claims alma (hmac/rsa)
		claims, err := validateToken(jwtSecretKey, tokenString)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// 5. UserID'yi güvenli bir şekilde alma (Type assertion check)
		userIdVal, ok := claims["sub"].(string)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token payload: sub missing or invalid")
		}

		// 6. Context'e güvenli key ile ekleme
		ctx = context.WithValue(ctx, UserIDKey, userIdVal)

		return handler(ctx, req)
	}
}

func validateToken(jwtSecretKey, tokenString string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Algoritma kontrolü (Önemli güvenlik adımı)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed // Unexpected signing method
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	// Claims type assertion kontrolü
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
