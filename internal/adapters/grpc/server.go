package grpc

import (
	"auth/internal/config"
	"auth/internal/ports"
	"auth_grpc"
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type GrpcServer struct {
	authS ports.AuthService
}

// NewGrpcServer returns gRPC server
func NewGrpcServer(authS ports.AuthService) *GrpcServer {
	return &GrpcServer{
		authS: authS,
	}
}

func (s *GrpcServer) Start(cfg *config.Configs) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Address, cfg.Grpc.Port)

	isReady := &atomic.Value{}
	isReady.Store(false)

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.Grpc.MaxConnectionIdle * time.Minute,
			MaxConnectionAge:  cfg.Grpc.MaxConnectionAge * time.Minute,
			Time:              cfg.Grpc.Timeout * time.Minute,
			Timeout:           cfg.Grpc.Timeout * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	as := NewAuthAPI(s.authS)

	auth_grpc.RegisterAuthServiceServer(grpcServer, as)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		log.Info().Msgf("GRPC server is listening: %s", grpcAddr)
		if err := grpcServer.Serve(l); err != nil {
			log.Fatal().Err(err).Msg("Failed running gRPC server")
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Info().Msg("The service is ready to accept requests")
	}()

	if cfg.App.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info().Msgf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Info().Msgf("ctx.Done: %v", done)
	}

	isReady.Store(false)

	grpcServer.GracefulStop()
	log.Info().Msgf("grpcServer shut down correctly")

	return nil
}
