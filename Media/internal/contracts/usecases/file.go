package usecases

import (
	"Media/internal/domain"
	"context"
	"github.com/google/uuid"
	"io"
)

type CreateFileDTO struct {
	Name     string
	AuthorID uuid.UUID
	Size     int64
	Content  io.Reader
}

type FileUseCaseInterface interface {
	CreateFile(ctx context.Context, dto CreateFileDTO) (uuid.UUID, error)
	GetFile(ctx context.Context, id uuid.UUID) (domain.File, error)
	DeleteFile(ctx context.Context, id uuid.UUID) error
}
