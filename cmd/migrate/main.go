package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/logger"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner_rotator/server_config.toml", "Path to configuration file")
}

type (
	StorageCtrl interface {
		Connect(ctx context.Context) error
		Migrate(ctx context.Context, migrate string) (err error)
		Close(ctx context.Context) error
	}
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewMigratorConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.Output)
	defer logg.Close()

	logg.Info(fmt.Sprintf("Connecting to %v:%v", config.Storage.Host, config.Storage.Port))

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
		logg.Error("failed to connect to storage: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	logg.Info("Appying migrations...")
	err = storage.Migrate(ctx, "/migrations")
	if err != nil {
		logg.Error(fmt.Sprintf("Failed to apply migrations: %v", err))
		os.Exit(1)
	}
	logg.Info("Done")
}
