package usecases

import (
	"Media/internal/contracts/usecases"
	"Media/internal/domain"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type FileRepositoryInterface interface {
	CreateFile(ctx context.Context, file domain.File) error
	GetFile(ctx context.Context, id uuid.UUID) (domain.File, error)
	DeleteFile(ctx context.Context, id uuid.UUID) error
}

var _ usecases.FileUseCaseInterface = &FileUseCase{}

type FileUseCase struct {
	Repository FileRepositoryInterface
	logger     *slog.Logger
}

func (f *FileUseCase) CreateFile(ctx context.Context, dto usecases.CreateFileDTO) (uuid.UUID, error) {
	id := uuid.New()
	content := make([]byte, dto.Size)
	_, err := dto.Content.Read(content)
	if err != nil {
		return uuid.Nil, err
	}
	file := domain.File{
		ID:       id,
		Name:     dto.Name,
		Size:     dto.Size,
		Content:  content,
		AuthorID: dto.AuthorID,
	}
	err = f.Repository.CreateFile(ctx, file)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (f *FileUseCase) GetFile(ctx context.Context, id uuid.UUID) (domain.File, error) {
	return f.Repository.GetFile(ctx, id)
}

func (f *FileUseCase) DeleteFile(ctx context.Context, id uuid.UUID) error {
	return f.Repository.DeleteFile(ctx, id)
}

func NewFileUseCase(repository FileRepositoryInterface, logger *slog.Logger) *FileUseCase {
	return &FileUseCase{
		Repository: repository,
		logger:     logger,
	}
}
