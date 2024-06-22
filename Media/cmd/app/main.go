package main

import (
	"Media/config"
	"Media/internal/infrastructure/repositories/minio"
	"Media/internal/infrastructure/server"
	"Media/internal/usecases"

	"context"
	"log/slog"
	"os"

	defaultLogger "log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Environment constants
const (
	Development = "dev"
	Production  = "prod"
	Test        = "test"
)

func startApp() error {
	// Read the configuration
	cfgPath := config.FetchPath()
	defaultLogger.Println("Using config path: ", cfgPath)
	cfg := config.MustParseConfig(cfgPath)

	// Create a new logger
	log := InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.Any("env", cfg.Env))

	// Create a new MinIO client
	mc, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: *cfg.Minio.UseSSL,
	})
	if err != nil {
		log.Error("Failed to create MinIO client", slog.Any("error", err.Error()))
		return err
	}

	//Check if the bucket exists
	if exists, err := mc.BucketExists(context.Background(), cfg.Minio.Bucket); err != nil {
		log.Error("Failed to check if bucket exists", slog.Any("error", err.Error()))
		return err
	} else if !exists {
		if err := mc.MakeBucket(context.Background(), cfg.Minio.Bucket, minio.MakeBucketOptions{}); err != nil {
			log.Error("Failed to create bucket", slog.Any("error", err.Error()))
			return err
		}
	} else {
		log.Info("Bucket already exists", slog.Any("bucket", cfg.Minio.Bucket))
	}

	// Create a new repository
	fileRepo := minioRepo.NewFileRepository(mc, cfg.Minio.Bucket, log)

	// Create a new use case
	uc := usecases.NewFileUseCase(fileRepo, log)

	// Create a new server
	srv := server.NewServer(cfg.Server.Address, uc, log)

	// Start the server
	if err := srv.Start(); err != nil {
		log.Error("Failed to start server", slog.Any("error", err.Error()))
		return err
	}

	return nil
}

func main() {
	if err := startApp(); err != nil {
		panic(err)
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
