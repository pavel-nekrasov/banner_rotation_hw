package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	rotatorapp "github.com/pavel-nekrasov/banner_rotation_hw/internal/app"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/logger"
	internalgrpc "github.com/pavel-nekrasov/banner_rotation_hw/internal/server/grpc"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner_rotator/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewRotatorConfig(configFile)
	log := logger.New(config.Logger.Level, config.Logger.Output)
	defer log.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := storage.New(
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.DBName,
		config.Storage.User,
		config.Storage.Password,
	)
	err := storage.Connect(ctx)
	if err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		log.Close()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	app := rotatorapp.New(log, storage, config.Cache)
	server := internalgrpc.NewServer(config.Endpoint.Host,
		config.Endpoint.GRPCPort,
		config.Endpoint.HTTPPort,
		log,
		app,
	)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop Banner Rotator server: " + err.Error())
		}
		log.Info("Shutted down")
	}()

	log.Info("Banner Rotator is running...")
	if err := server.Start(ctx); err != nil {
		log.Error("failed to start rotator server: " + err.Error())
		cancel()
		os.Exit(1)
	}
	log.Info("Banner Rotator is stopped")
}
