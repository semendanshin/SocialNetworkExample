package usecases

import (
	"Posts/internal/domain"
	"context"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=PostUseCase

// PostUseCase is a use case for posts.
type PostUseCase interface {
	AbstractUseCaseInterface[*domain.Post]
	GetByAuthorID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*domain.Post, error)
}
