package main

import (
	"codegen/internal/config"
	"codegen/internal/repository/pg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpctransport "codegen/internal/transport/grpc"
)

func main() {
	cfg := config.LoadFromEnv()
	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DBUrl, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := pg.NewPGRepo(pool)
	// construct services...
	s := grpctransport.NewServer(nil, nil) // pass services

	go func() {
		if err := s.Serve(cfg.Port); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
	defer cancel()
	s.GracefulStop(timeoutCtx)
}
