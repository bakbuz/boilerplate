package main

import (
	"codegen/internal/config"
	"codegen/internal/database"
	"codegen/internal/interceptor"
	"codegen/internal/repository"
	"codegen/internal/service"
	"codegen/proto/pb"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Config Yükle
	cfg := config.Load()

	// 2. DB Bağlantısı
	ctx := context.Background()
	dbPool, err := database.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer dbPool.Close()

	// 3. Repository & Services
	repo := repository.NewPostgresRepo(dbPool)
	productSvc := service.NewProductServer(repo)

	// 4. gRPC Server Setup
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Optimizasyon: KeepAlive parametreleri
	// opts := []grpc.ServerOption{
	// 	grpc.KeepaliveParams(keepalive.ServerParameters{...}),
	//  grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.JWTSecret)),
	// }

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.JWTSecret)),
	)

	pb.RegisterProductServiceServer(s, productSvc)
	// pb.RegisterOrderServiceServer(s, orderSvc)

	// Development modunda Reflection aç (Postman/gRPCurl için)
	if cfg.AppEnv != "production" {
		reflection.Register(s)
	}

	// 5. Graceful Shutdown (Systemd için kritik)
	go func() {
		log.Printf("Server listening on %s", cfg.Port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
