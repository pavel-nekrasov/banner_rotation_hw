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

	dbConn := storage.NewConnection(config.Storage)

	logg.Info("Applying migrations...")
	err := dbConn.Migrate(ctx, "/migrations")
	if err != nil {
		logg.Error(fmt.Sprintf("Failed to apply migrations: %v", err))
		cancel()
		os.Exit(1) //nolint: gocritic
	}
	logg.Info("Done")
}
