package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/evg555/antibrutforce/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	srv    *grpc.Server
	logger Logger
	app    Application
	cfg    *config.Config
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface{}

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
		cfg:    &cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.loggingMiddleware),
	)

	//reflection.Register(s.srv)
	//pb.RegisterEventServiceServer(s.srv, Handler{
	//	app:    s.app,
	//	logger: s.logger,
	//})

	addr := net.JoinHostPort(s.cfg.App.Host, s.cfg.App.Port)

	//listener, err := net.Listen("tcp", addr)
	//if err != nil {
	//	return err
	//}

	s.logger.Info(fmt.Sprintf("grpc server starting at %s", addr))

	//if err = s.srv.Serve(listener); err != nil {
	//	return err
	//}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.logger.Info("grpc server stopping...")
	s.srv.Stop()
	return nil
}
