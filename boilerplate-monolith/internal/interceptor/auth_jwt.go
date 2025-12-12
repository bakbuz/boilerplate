package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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

	token := strings.TrimPrefix(tokens[0], "Bearer ")
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid token format")
	}

	// Validate token
	claims, err := ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Add claims to context
	ctx = context.WithValue(ctx, "user_id", claims.UserID)

	return handler(ctx, req)
}

func isPublicEndpoint(method string) bool {
	publicEndpoints := []string{
		"/api.v1.AuthService/Login",
		"/api.v1.AuthService/Register",
		"/grpc.health.v1.Health/Check",
	}

	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}

	return false
}

/*

func ValidateToken(tokenStr, secret string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(secret), nil
	})
	if err != nil || !t.Valid {
		return "", err
	}
	claims := t.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	return sub, nil
}


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
*/
