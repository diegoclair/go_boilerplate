package shutdown

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

const gracefulShutdownTimeout = 10 * time.Second

type ShutdownOptions func(s *shutdown)

type shutdown struct {
	restServer *echo.Echo
	grpcServer *grpc.Server
	listener   net.Listener
}

func GracefulShutdown(ctx context.Context, log logger.Logger, opts ...ShutdownOptions) {
	s := &shutdown{}

	for _, opt := range opts {
		opt(s)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	<-stop

	log.Info(ctx, "Shutting down server...")

	if s.restServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		err := s.restServer.Shutdown(ctx)
		if err != nil {
			log.Errorw(ctx, "Failed to shutdown rest server", logger.Err(err))
		}
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.listener != nil {
		s.listener.Close()
	}
}

func WithRestServer(restServer *echo.Echo) ShutdownOptions {
	return func(s *shutdown) {
		s.restServer = restServer
	}
}

func WithGrpcServer(grpcServer *grpc.Server) ShutdownOptions {
	return func(s *shutdown) {
		s.grpcServer = grpcServer
	}
}

func WithListener(listener net.Listener) ShutdownOptions {
	return func(s *shutdown) {
		s.listener = listener
	}
}
