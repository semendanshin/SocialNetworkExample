package graph

import (
	"Posts/internal/infrastructure/graph/middleware"
	"Posts/internal/interfaces/usecases"
	"Posts/pkg/jwtservice"
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const defaultPort = "8080"

// Server is a GraphQL server.
type Server struct {
	port             string
	jwtGen           *jwtservice.Service
	logger           *slog.Logger
	schema           graphql.ExecutableSchema
	srv              http.Server
	enablePlayground bool

	postUseCase    usecases.PostUseCase
	commentUseCase usecases.CommentUseCase
	userUseCase    usecases.UserUseCase
}

// NewServer creates a new server.
func NewServer(
	port string,
	jwtGen *jwtservice.Service,
	logger *slog.Logger,
	schema graphql.ExecutableSchema,
	enablePlayground bool,
	postUseCase usecases.PostUseCase,
	commentUseCase usecases.CommentUseCase,
	userUseCase usecases.UserUseCase,
) *Server {
	return &Server{
		port:             port,
		jwtGen:           jwtGen,
		logger:           logger,
		schema:           schema,
		enablePlayground: enablePlayground,
		postUseCase:      postUseCase,
		commentUseCase:   commentUseCase,
		userUseCase:      userUseCase,
	}
}

// Run starts the server.
func (s *Server) Run() error {
	if s.port = os.Getenv("PORT"); s.port == "" {
		s.port = defaultPort
	}

	router := mux.NewRouter()

	if s.enablePlayground {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	graphQlHandler := handler.NewDefaultServer(s.schema)
	graphQlHandler.Use(extension.FixedComplexityLimit(1000))

	queryRouter := router.PathPrefix("/query").Subrouter()
	queryRouter.Use(middleware.DataLoader(s.postUseCase, s.commentUseCase, s.userUseCase, s.logger))
	queryRouter.Use(middleware.Auth(s.jwtGen, s.logger))
	queryRouter.Handle("", graphQlHandler)

	s.logger.Info("starting server", slog.Any("port", s.port))

	s.srv = http.Server{
		Addr:    ":" + s.port,
		Handler: router,
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		s.logger.Info("gratefully shutting down server")
		if err := s.srv.Shutdown(context.Background()); err != nil {
			s.logger.Error("failed to shutdown server", slog.Any("error", err.Error()))
		}
	}()

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to listen and serve", slog.Any("error", err.Error()))
	}

	s.logger.Info("server stopped")

	return nil
}
