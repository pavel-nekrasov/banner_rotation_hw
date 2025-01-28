package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	rotatorapp "github.com/pavel-nekrasov/banner_rotation_hw/internal/app"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/logger"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/queue"
	internalgrpc "github.com/pavel-nekrasov/banner_rotation_hw/internal/server/grpc"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner_rotator/server_config.toml", "Path to configuration file")
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

	// db
	dbConn := storage.NewConnection(
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.DBName,
		config.Storage.User,
		config.Storage.Password,
	)

	storage := storage.NewStorage(
		dbConn,
	)
	err := dbConn.Connect(ctx)
	if err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		log.Close()
		os.Exit(1) //nolint:gocritic
	}
	defer dbConn.Close()

	// queue
	queueConn := queue.NewConnection(config.Queue.QueueServerConf)
	if err := queueConn.Connect(); err != nil {
		log.Error(fmt.Sprintf("failed to connect to queue: %s", err.Error()))
		dbConn.Close()
		log.Close()
		os.Exit(1)
	}
	defer queueConn.Close()

	queueProducer := queue.NewProducer(queueConn, config.Queue)
	if err := queueProducer.Start(); err != nil {
		log.Error(fmt.Sprintf("failed to create exchange: %s", err.Error()))
		queueConn.Close()
		dbConn.Close()
		log.Close()
		os.Exit(1)
	}
	defer queueProducer.Close()

	// grpc service
	app := rotatorapp.New(ctx, log, storage, queueProducer, config.Cache)
	server := internalgrpc.NewServer(config.Endpoint.Host,
		config.Endpoint.GRPCPort,
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
		queueProducer.Close()
		queueConn.Close()
		dbConn.Close()
		log.Close()
		os.Exit(1)
	}
	log.Info("Banner Rotator is stopped")
}
