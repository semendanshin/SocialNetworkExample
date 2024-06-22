package server

import (
	"Media/internal/infrastructure/server/handlers"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"os"
	"os/signal"
	"syscall"

	"Media/internal/contracts/usecases"

	"log/slog"
	"net/http"
)

type Server struct {
	address string
	fuc     usecases.FileUseCaseInterface
	logger  *slog.Logger
	server  *http.Server
}

func NewServer(address string, fuc usecases.FileUseCaseInterface, logger *slog.Logger) *Server {
	return &Server{
		address: address,
		fuc:     fuc,
		logger:  logger,
	}
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	fileHandler := handlers.NewFileHandler(s.fuc, s.logger)
	fileHandler.RegisterRoutes(router)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: router,
	}

	s.logger.Info("starting server", slog.Any("address", s.address))
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		s.logger.Info("gratefully shutting down server")
		if err := s.server.Shutdown(context.Background()); err != nil {
			s.logger.Error("failed to shutdown server", slog.Any("error", err.Error()))
		}
	}()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to listen and serve", slog.Any("error", err.Error()))
	}

	s.logger.Info("server stopped")

	return nil
}
