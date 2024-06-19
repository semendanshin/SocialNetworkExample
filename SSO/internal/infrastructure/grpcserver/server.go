package grpcserver

import (
	pb "SSO/gen/go"
	"SSO/internal/contracts/usecases"
	"SSO/internal/infrastructure/grpcserver/interceptors"
	"SSO/internal/infrastructure/grpcserver/services"
	"SSO/pkg/jwt"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Server is the interface that wraps the Serve method.
type Server interface {
	Serve() error
}

type server struct {
	logger           *slog.Logger
	grpcServer       *grpc.Server
	address          string
	enableReflection bool

	jwtManager *jwt.Manager
	uuc        usecases.UserUseCases
}

// NewServer returns a new instance of the Server.
func NewServer(logger *slog.Logger, address string, enableReflection bool, jwtManager *jwt.Manager, uuc usecases.UserUseCases) Server {
	return &server{
		logger:           logger,
		address:          address,
		enableReflection: enableReflection,
		jwtManager:       jwtManager,
		uuc:              uuc,
	}
}

// Serve starts the server.
func (s *server) Serve() error {

	lis, err := net.Listen("tcp4", s.address)
	if err != nil {
		s.logger.Error("failed to listen: %v", err)
		return err
	}

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthInterceptor(s.jwtManager, s.logger),
		),
	)

	pb.RegisterUserServiceServer(s.grpcServer, services.NewUserServiceServer(s.logger, s.uuc))
	pb.RegisterAuthServiceServer(s.grpcServer, services.NewAuthServiceServer(s.logger, s.uuc, s.jwtManager))

	if s.enableReflection {
		reflection.Register(s.grpcServer)
	}

	s.logger.Info("Starting server")
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		s.grpcServer.Stop()
		s.logger.Info("Gracefully stopping server")
	}()

	if err := s.grpcServer.Serve(lis); err != nil {
		if !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error("Failed to start server", slog.String("error", err.Error()))
			return err
		}
	}

	s.logger.Info("Server stopped")

	return nil
}
