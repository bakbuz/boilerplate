package auth

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