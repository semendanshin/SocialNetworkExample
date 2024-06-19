package usecases

import (
	"Posts/internal/domain"
	usecaseInterfaces "Posts/internal/interfaces/usecases"
	"context"
	"github.com/google/uuid"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=CommentRepository

// CommentRepository is a repository for comments.
type CommentRepository interface {
	usecaseInterfaces.AbstractRepositoryInterface[*domain.Comment]
	GetChildren(ctx context.Context, commentID uuid.UUID, limit int, offset int) ([]*domain.Comment, error)
	GetByPostID(ctx context.Context, postID uuid.UUID, limit int, offset int) ([]*domain.Comment, error)
	GetLastComment(ctx context.Context, postID uuid.UUID, lastSeen time.Time, limit int) ([]*domain.Comment, error)
}

var _ usecaseInterfaces.CommentUseCase = &CommentUseCase{}

// CommentUseCase is a use case for comments.
type CommentUseCase struct {
	Repository CommentRepository
	usecaseInterfaces.AbstractUseCase[*domain.Comment]
}

// NewCommentUseCase creates a new CommentUseCase.
func NewCommentUseCase(repository CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		Repository:      repository,
		AbstractUseCase: usecaseInterfaces.NewAbstractUseCase[*domain.Comment](repository),
	}
}

// GetChildren returns all children of a comment.
func (uc *CommentUseCase) GetChildren(ctx context.Context, commentID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	return uc.Repository.GetChildren(ctx, commentID, limit, offset)
}

// GetByPostID returns all comments of a post.
func (uc *CommentUseCase) GetByPostID(ctx context.Context, postID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	return uc.Repository.GetByPostID(ctx, postID, limit, offset)
}

// Create creates a new comment.
func (uc *CommentUseCase) Create(ctx context.Context, entity *domain.Comment) error {
	if len(entity.Content) > 2000 {
		return domain.ErrCommentIsTooLong
	}
	entity.SetID(uuid.New())
	return uc.AbstractUseCase.Create(ctx, entity)
}

// GetLastComment returns the last comments of a post.
func (uc *CommentUseCase) GetLastComment(ctx context.Context, postID uuid.UUID, lastSeen time.Time, limit int) ([]*domain.Comment, error) {
	return uc.Repository.GetLastComment(ctx, postID, lastSeen, limit)
}
