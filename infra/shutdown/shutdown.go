package shutdown

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegoclair/logger"
	"google.golang.org/grpc"
)

type Option func(s *stopper)

type stopper struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// Context returns a context cancelled when the process is asked to terminate.
// The REST server drains on that cancellation, so whatever must outlive the
// requests in flight is stopped by Stop, after serving returns.
func Context(ctx context.Context, log logger.Logger) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		os.Interrupt,
	)

	go func() {
		<-signals
		log.Info(ctx, "Shutting down server...")
		cancel()
	}()

	return ctx
}

func Stop(ctx context.Context, log logger.Logger, opts ...Option) {
	s := &stopper{}
	for _, opt := range opts {
		opt(s)
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Error(ctx, "Failed to close the listener", logger.Err(err))
		}
	}
}

func WithGrpcServer(grpcServer *grpc.Server) Option {
	return func(s *stopper) {
		s.grpcServer = grpcServer
	}
}

func WithListener(listener net.Listener) Option {
	return func(s *stopper) {
		s.listener = listener
	}
}
