package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/evg555/antibrutforce/api/pb"
	"github.com/evg555/antibrutforce/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

type Application interface {
	AddIPWhitelist(ctx context.Context, subnet string) error
	DeleteIPWhitelist(ctx context.Context, subnet string) error
	AddIPBlacklist(ctx context.Context, subnet string) error
	DeleteIPBlacklist(ctx context.Context, subnet string) error

	IsInBlacklist(ctx context.Context, ip string) bool
	IsInWhitelist(ctx context.Context, ip string) bool
	HasLimits(login, password, ip string) bool

	ResetBucket(password, ip string)
}

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

	reflection.Register(s.srv)
	pb.RegisterAppServiceServer(s.srv, Handler{
		app:    s.app,
		logger: s.logger,
	})

	addr := net.JoinHostPort(s.cfg.App.Host, s.cfg.App.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("grpc server starting at %s", addr))

	if err = s.srv.Serve(listener); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.logger.Info("grpc server stopping...")
	s.srv.Stop()
	return nil
}
