package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evg555/antibrutforce/internal/app"
	"github.com/evg555/antibrutforce/internal/config"
	"github.com/evg555/antibrutforce/internal/logger"
	"github.com/evg555/antibrutforce/internal/ratelimiter"
	internalgrpc "github.com/evg555/antibrutforce/internal/server/grpc"
	sqlstorage "github.com/evg555/antibrutforce/internal/storage/sql"
)

func main() {
	flag.Parse()

	exitCode := 0

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cfg := config.NewConfig()
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Format)

	storage := sqlstorage.New(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.DBName,
	)

	rateLimiter := ratelimiter.NewAuthRateLimiter(ctx, cfg.RateLimiter)
	service := app.New(&logg, storage, rateLimiter)
	server := internalgrpc.NewServer(cfg, &logg, service)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := storage.Close(ctx); err != nil {
			logg.Error("failed to close connection to storage: " + err.Error())
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop server: " + err.Error())
		}

		os.Exit(exitCode)
	}()

	logg.Info("app is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start server: " + err.Error())

		exitCode = 1
		cancel()
	}

	select {}
}
