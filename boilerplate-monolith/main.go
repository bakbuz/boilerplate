package main

import (
	"codegen/api/pb"
	"codegen/internal/bootstrap"
	"codegen/internal/database"
	"codegen/internal/handler"
	"codegen/internal/interceptor"
	"codegen/internal/repository"
	"codegen/internal/service"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	// Zerolog stack trace ayarları
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func main() {
	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer quit()

	// Log çıktısını ayarla
	logWriter := zerolog.SyncWriter(os.Stdout)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	logger := zerolog.New(logWriter).
		With().
		Timestamp().
		Caller().
		Stack().
		Logger()

	ctx = logger.WithContext(ctx)

	flag.Parse()

	if err := run(ctx, &logger); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

func run(ctx context.Context, logger *zerolog.Logger) error {
	env := strings.TrimSpace(os.Getenv("ENV"))
	if env == "" {
		env = "dev"
	}

	configFile := fmt.Sprintf("config.%s.json", env)
	logger.Debug().Str("ENV", env).Str("FILE", configFile).Msg("")

	cfg, err := bootstrap.LoadConfig(configFile)
	if err != nil {
		return errors.WithMessage(err, "failed to read configuration file")
	}

	// database
	dbPool, err := database.NewPool(ctx, cfg.DataSources.Default)
	if err != nil {
		return errors.WithMessage(err, "failed to connect the database")
	}
	defer dbPool.Close()
	logger.Info().Interface("Version", dbPool.PostgresVersion(ctx)).Msg("database connected")

	// repositories
	brandRepo := repository.NewBrandRepository(dbPool)
	productRepo := repository.NewProductRepository(dbPool)

	// services
	brandSvc := service.NewBrandService(brandRepo)
	productSvc := service.NewProductService(productRepo)

	// handlers
	demoHandler := handler.NewDemoHandler()
	brandHandler := handler.NewBrandHandler(brandSvc)
	productHandler := handler.NewProductHandler(productSvc)

	// gRPC server options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 20), // 20 MB max message size
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.Jwt.SecretKey)),
	}

	// gRPC server instance
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDemoServiceServer(grpcServer, demoHandler)
	pb.RegisterBrandServiceServer(grpcServer, brandHandler)
	pb.RegisterCatalogServiceServer(grpcServer, productHandler)

	// Development modunda Reflection aç (Postman/gRPCurl için)
	if env != "production" {
		reflection.Register(grpcServer)
	}

	// gRPC server listen
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Panic().Err(err).Msgf("failed to listen on port: 50051")
	}

	logger.Info().Msgf("server listening at 50051")

	// gRPC server start
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error().Msgf("failed to serve: %v", err)
		}
	}()

	// Sistem sinyallerini dinle
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan // SIGINT veya SIGTERM sinyali alınana kadar bekle
	logger.Info().Msg("Shutting down gRPC server gracefully...")

	// Sunucuyu durdur
	gracefulStop(grpcServer, logger)

	logger.Info().Msg("Server stopped.")

	return nil
}

// Sunucuyu düzgün bir şekilde kapatır
func gracefulStop(grpcServer *grpc.Server, logger *zerolog.Logger) {
	// Bağlantıları tamamlaması için bir süre tanır
	const shutdownTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// gRPC sunucusunu kibarca durdur
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Info().Msg("The server has been successfully shut down.") // Sunucu başarılı bir şekilde kapatıldı.
	case <-ctx.Done():
		logger.Info().Msg("Shutdown timed out; the server is being forcibly stopped.") // Kapatma zaman aşımına uğradı, sunucu zorla durduruluyor.
		grpcServer.Stop()                                                              // Zorla durdur
	}
}
