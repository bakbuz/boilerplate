package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"codegen/internal/bootstrap"
	"codegen/internal/database"
	"codegen/internal/handler"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// @title					Codegen API
// @version				1.0
// @description		Codegen service.
// @contact.name	Bayram Akbuz
// @contact.url		https://github.com/bakbuz
// @BasePath			/
func main() {
	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt)
	defer quit()

	logWriter := zerolog.SyncWriter(os.Stdout)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	ctx = logger.WithContext(ctx)

	flag.Parse()

	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

func run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	logger.Debug().Str("ENV", env).Msg("")

	conf, err := bootstrap.NewConfig(fmt.Sprintf("config.%s.json", env))
	if err != nil {
		return errors.WithMessage(err, "failed to read configuration file")
	}
	//logger.Debug().Any("conf", conf).Msg("")

	// database
	db, err := database.New(conf.DataSources.Default)
	if err != nil {
		return errors.WithMessage(err, "failed to connect the database")
	}
	logger.Info().Msg("database connected")

	// echo
	e := echo.New()
	e.Use(middleware.BodyLimit("2M"))
	e.Pre(middleware.MethodOverride())
	e.Validator = handler.NewCustomValidator()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{},
		AllowCredentials: true,
	}))

	// bootstraper
	bootstrap.RegisterRoutes(e, db, conf)

	// server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.App.Port),
		Handler: e,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	logger.Info().Int("port", conf.App.Port).Msg("application is starting")

	srvErrCh := make(chan error)
	go func() {
		defer close(srvErrCh)

		if err := srv.ListenAndServe(); err != nil {
			srvErrCh <- errors.WithStack(err)
		}

	}()

	select {
	case err := <-srvErrCh:
		return errors.WithStack(err)

	case <-ctx.Done():
		logger.Debug().Msg("graceful shutdown has been started.")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return errors.WithStack(err)
		}

		logger.Debug().Msg("graceful shutdown has been completed.")

		return nil
	}
}
