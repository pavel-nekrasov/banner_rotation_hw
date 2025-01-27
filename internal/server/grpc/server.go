package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/server/grpc/pb"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/server/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	host     string
	grpcPort int
	httpPort int
	logger   Logger
	app      services.Application
	server   *grpc.Server
	gwServer *http.Server
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

func NewServer(host string, grpcPort int, httpPort int, logger Logger, app services.Application) *Server {
	return &Server{
		host:     host,
		grpcPort: grpcPort,
		httpPort: httpPort,
		logger:   logger,
		app:      app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	grpcBindAddr := fmt.Sprintf("%v:%v", s.host, s.grpcPort)
	s.logger.Info(fmt.Sprintf("Starting GRPC on %v...", grpcBindAddr))
	lsn, err := net.Listen("tcp", grpcBindAddr)
	if err != nil {
		return err
	}
	/*
		httpBindAddr := fmt.Sprintf("%v:%v", s.host, s.httpPort)
		s.logger.Info(fmt.Sprintf("Starting HTTP on %v...", httpBindAddr))

			gwClient, err := grpc.NewClient(grpcBindAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}

			gwMux := runtime.NewServeMux()
			err = pb.RegisterEventsHandler(ctx, gwMux, gwClient)
			if err != nil {
				return err
			}

			gwMuxWithLogging := httpmiddleware.NewLoggingMiddleware(s.logger, gwMux)
			s.gwServer = &http.Server{Addr: httpBindAddr, Handler: gwMuxWithLogging, ReadTimeout: time.Second * 10}
			go func() {
				if err := s.gwServer.ListenAndServe(); err != http.ErrServerClosed {
					panic(err)
				}
			}()
	*/

	s.server = grpc.NewServer()
	pb.RegisterBannersServer(s.server, services.NewBannerService(s.logger, s.app))
	reflection.Register(s.server)

	go func() {
		if err := s.server.Serve(lsn); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping adapters...")

	defer s.server.Stop()
	return s.gwServer.Shutdown(ctx)
}
