package main

import (
	"Posts/internal/infrastructure/graph"
	"Posts/internal/infrastructure/graph/resolvers"
	"Posts/internal/infrastructure/repository/sql"
	"Posts/internal/usecases"
	"Posts/pkg/jwtservice"
	"fmt"
	"gorm.io/driver/postgres"

	inmemory "Posts/internal/infrastructure/repository/in-memory"
	defaultLogger "log"

	"gorm.io/gorm"

	"Posts/config"
	"log/slog"
	"os"
)

// Env constants
const (
	Test        = "test"
	Production  = "prod"
	Development = "dev"
)

func main() {
	// Read config
	cfgPath := config.FetchPath()
	defaultLogger.Println("Using config path: ", cfgPath)
	cfg := config.MustParseConfig(cfgPath)

	// Init logger
	log := InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.Any("env", cfg.Env))

	var db *gorm.DB
	var err error

	// Init database
	if cfg.UseDatabase == nil || !*cfg.UseDatabase {
		log.Info("Using in-memory database")
	} else {
		log.Info("Using Postgres database", slog.Any("host", cfg.Postgres.Host))

		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.Name)

		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{TranslateError: true})
		if err != nil {
			log.Error("Failed to connect to database", slog.Any("error", err.Error()))
			return
		}

		// Ping database
		if err = db.Exec("SELECT 1").Error; err != nil {
			log.Error("Failed to ping database", slog.Any("error", err.Error()))
			return
		}
	}

	// Init repositories
	var postRepo usecases.PostRepository
	var commentRepo usecases.CommentRepository
	var userRepo usecases.UserRepository

	if cfg.UseDatabase == nil || !*cfg.UseDatabase {
		postRepo = inmemory.NewPostInMemoryRepository(log)
		commentRepo = inmemory.NewCommentInMemoryRepository(log)
		userRepo = inmemory.NewUserInMemoryRepository(log)
	} else {
		postRepo = sql.NewPostSQLRepository(db, log)
		commentRepo = sql.NewCommentSQLRepository(db, log)
		userRepo = sql.NewUserSQLRepository(db, log)
	}

	// Init UseCases
	postUseCase := usecases.NewPostUseCase(postRepo)
	commentUseCase := usecases.NewCommentUseCase(commentRepo)
	userUseCase := usecases.NewUserUseCase(userRepo)

	// Init Resolver and Schema
	resolver := resolvers.NewResolver(
		postUseCase,
		commentUseCase,
		userUseCase,
		log,
	)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})

	// Init jwt Service
	jwtGen := jwtservice.NewGenerator(cfg.Tokens.Secret, cfg.Tokens.AccessTTL, cfg.Tokens.RefreshTTL)

	// Init server
	srv := graph.NewServer(
		"8080",
		jwtGen,
		log,
		schema,
		true,
		postUseCase,
		commentUseCase,
		userUseCase,
	)

	// Run server
	if err := srv.Run(); err != nil {
		log.Error("Failed to run server", slog.Any("error", err.Error()))
		os.Exit(1)
	}
}

// InitLogger initializes a logger based on the environment.
func InitLogger(env string) *slog.Logger {
	var log *slog.Logger
	if env == Production {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		return log
	} else if env == Test {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		return log
	} else {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		return log
	}
}
