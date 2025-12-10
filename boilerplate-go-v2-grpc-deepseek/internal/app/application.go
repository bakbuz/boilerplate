package app

import (
    "context"
    "net"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive"

    "github.com/yourusername/grpc-highperf-backend/configs"
    "github.com/yourusername/grpc-highperf-backend/internal/handler"
    "github.com/yourusername/grpc-highperf-backend/internal/pkg/database"
    "github.com/yourusername/grpc-highperf-backend/internal/repository/postgres"
    "github.com/yourusername/grpc-highperf-backend/internal/service"
)

type Application struct {
    config     *configs.Config
    logger     zerolog.Logger
    grpcServer *grpc.Server
    dbPool     *pgxpool.Pool
}

func NewApplication(ctx context.Context, cfg *configs.Config, log zerolog.Logger) (*Application, error) {
    // Database connection with pooling and retry
    dbPool, err := database.NewPostgresPool(ctx, cfg.Database)
    if err != nil {
        return nil, err
    }

    // Setup gRPC server with performance optimizations
    grpcServer := grpc.NewServer(
        grpc.KeepaliveParams(keepalive.ServerParameters{
            MaxConnectionIdle: cfg.Server.GRPCKeepAlive,
            Time:              cfg.Server.GRPCKeepAliveTime,
            Timeout:           cfg.Server.GRPCKeepAliveTimeout,
        }),
        grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
            MinTime:             cfg.Server.GRPCKeepAliveTime / 2,
            PermitWithoutStream: true,
        }),
        grpc.MaxRecvMsgSize(cfg.Server.MaxMsgSize),
        grpc.MaxSendMsgSize(cfg.Server.MaxMsgSize),
        grpc.NumStreamWorkers(uint32(cfg.Server.WorkerCount)),
    )

    // Initialize repositories
    productRepo := postgres.NewProductRepository(dbPool)
    orderRepo := postgres.NewOrderRepository(dbPool)

    // Initialize services
    productService := service.NewProductService(productRepo)
    orderService := service.NewOrderService(orderRepo)

    // Register handlers
    handler.RegisterProductHandler(grpcServer, productService)
    handler.RegisterOrderHandler(grpcServer, orderService)

    return &Application{
        config:     cfg,
        logger:     log,
        grpcServer: grpcServer,
        dbPool:     dbPool,
    }, nil
}

func (a *Application) Start(ctx context.Context) error {
    lis, err := net.Listen("tcp", a.config.Server.Address)
    if err != nil {
        return err
    }

    a.logger.Info().Str("address", a.config.Server.Address).Msg("Starting gRPC server")
    return a.grpcServer.Serve(lis)
}

func (a *Application) Shutdown(ctx context.Context) error {
    a.grpcServer.GracefulStop()
    a.dbPool.Close()
    return nil
}