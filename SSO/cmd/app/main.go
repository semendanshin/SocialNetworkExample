package main

import (
	"SSO/config"
	"SSO/internal/infrastructure/grpcserver"
	gormrepository "SSO/internal/infrastructure/repositories/gorm"
	"SSO/internal/usecases"
	"SSO/pkg/jwt"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	defaultLogger "log"
	"log/slog"
	"os"
)

const (
	Development = "dev"
	Production  = "prod"
	Test        = "test"
)

func startApp() error {
	// Load configuration
	path := config.FetchPath()
	defaultLogger.Println("Using config path: ", path)
	cfg := config.MustParseConfig(path)

	// Initialize logger
	log := InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.Any("env", cfg.Env))

	// Initialize database
	log.Info("Using Postgres database", slog.Any("host", cfg.Postgres.Host))

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.Name)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Error("Failed to connect to database", slog.Any("error", err.Error()))
		return err
	}

	// Ping database
	if err = db.Exec("SELECT 1").Error; err != nil {
		log.Error("Failed to ping database", slog.Any("error", err.Error()))
		return err
	}

	// Initialize repositories
	userRepo := gormrepository.NewGormUserRepository(
		db,
		log,
	)

	// Initialize useCases
	userUseCases := usecases.NewUserUseCases(
		userRepo,
	)

	// Initialize jwt manager
	publicKey, err := jwt.ReadPublicKey(cfg.Tokens.PublicKeyPath)
	if err != nil {
		log.Error("Failed to read public key", slog.Any("error", err.Error()))
		return err
	}

	privateKey, err := jwt.ReadPrivateKey(cfg.Tokens.PrivateKeyPath)
	if err != nil {
		log.Error("Failed to read private key", slog.Any("error", err.Error()))
		return err
	}

	managerOptions := jwt.ManagerOptions{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  cfg.Tokens.AccessTTL,
		RefreshTTL: cfg.Tokens.RefreshTTL,
	}

	jwtManager := jwt.NewManager(managerOptions)

	// Initialize server
	server := grpcserver.NewServer(
		log,
		cfg.Server.Address,
		true,
		jwtManager,
		userUseCases,
	)

	// Start server
	if err := server.Serve(); err != nil {
		log.Error("Failed to start server", slog.Any("error", err.Error()))
		return err
	}

	return nil
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

func main() {
	if err := startApp(); err != nil {
		panic(err)
	}
}
