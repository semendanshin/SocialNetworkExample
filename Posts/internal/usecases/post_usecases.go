package usecases

import (
	"Posts/internal/domain"
	usecaseInterfaces "Posts/internal/interfaces/usecases"
	"context"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=PostRepository

// PostRepository is a repository for posts.
type PostRepository interface {
	usecaseInterfaces.AbstractRepositoryInterface[*domain.Post]
	GetByAuthorID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*domain.Post, error)
}

var _ usecaseInterfaces.PostUseCase = &PostUseCase{}

// PostUseCase is a use case for posts.
type PostUseCase struct {
	Repository PostRepository
	usecaseInterfaces.AbstractUseCase[*domain.Post]
}

// NewPostUseCase creates a new PostUseCase.
func NewPostUseCase(repository PostRepository) *PostUseCase {
	return &PostUseCase{
		Repository:      repository,
		AbstractUseCase: usecaseInterfaces.NewAbstractUseCase[*domain.Post](repository),
	}
}

// GetByAuthorID returns all posts by a user.
func (uc *PostUseCase) GetByAuthorID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*domain.Post, error) {
	return uc.Repository.GetByAuthorID(ctx, userID, limit, offset)
}
