// Package app implements keeperd service.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/config"
	cgrpc "github.com/derpartizanen/gophkeeper/internal/keeperd/controller/grpc"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/grpcserver"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/logger"
)

const (
	MinimalSecretLength = 32
	ShutdownTimeout     = 60 * time.Second
)

// Run initializes and starts the keeperd service.
func Run(cfg *config.Config) error {
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("app - Run - logger.New: %w", err)
	}

	log.Info().Msg(cfg.String())

	if len([]byte(cfg.Secret)) < MinimalSecretLength {
		log.Warn().Msg("Insecure signature: secret key is shorter than 32 bytes!")
	}

	pg, err := postgres.New(string(cfg.DatabaseURI), log)
	if err != nil {
		return fmt.Errorf("app - Run - postgres.New: %w", err)
	}

	repos := repo.New(pg)
	services := service.New(cfg, repos)

	grpcSrv, err := grpcserver.New(
		cfg.Address,
		cfg.CrtPath,
		cfg.KeyPath,
		grpc.MaxRecvMsgSize(cgrpc.DefaultMaxMessageSize),
		grpc.ChainUnaryInterceptor(
			cgrpc.LoggingUnaryInterceptor(log),
			cgrpc.AuthUnaryInterceptor(cfg.Secret),
		),
	)
	if err != nil {
		return fmt.Errorf("app - Run - grpcserver.New: %w", err)
	}

	cgrpc.RegisterRoutes(grpcSrv.Instance(), services)
	grpcSrv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - interrupt: signal " + s.String())

	case err := <-grpcSrv.Notify():
		log.Error().Err(err).Msg("app - Run - grpcSrv.Notify")
	}

	log.Info().Msg("Shutting down...")

	stopped := make(chan struct{})

	stopCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	go func() {
		shutdown(log, grpcSrv, pg)
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info().Msg("Service shutdown successful")

	case <-stopCtx.Done():
		log.Warn().Msgf("Exceeded %s shutdown timeout, exit forcibly", ShutdownTimeout)
	}

	return nil
}

func shutdown(
	log *logger.Logger,
	grpcSrv *grpcserver.Server,
	pg *postgres.Postgres,
) {
	log.Info().Msg("Shutting down gRPC API...")
	grpcSrv.Shutdown()

	log.Info().Msg("Shutting down database connection...")
	pg.Close()
}
