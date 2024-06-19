package inmemory

import (
	"Posts/internal/domain"
	"Posts/internal/usecases"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

var _ usecases.CommentRepository = &CommentInMemoryRepository{}

// CommentInMemoryRepository is a repository for comments.
type CommentInMemoryRepository struct {
	AbstractInMemoryRepository[*domain.Comment]
}

// NewCommentInMemoryRepository creates a new CommentInMemoryRepository.
func NewCommentInMemoryRepository(logger *slog.Logger) *CommentInMemoryRepository {
	return &CommentInMemoryRepository{
		AbstractInMemoryRepository: NewAbstractInMemoryRepository[*domain.Comment](logger),
	}
}

// GetChildren returns all children of a comment.
func (r *CommentInMemoryRepository) GetChildren(ctx context.Context, commentID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var keys []uuid.UUID
	for key := range r.entities {
		keys = append(keys, key)
	}

	var comments []*domain.Comment
	i := offset
	for i < len(keys) && len(comments) < limit {
		if r.entities[keys[i]].ParentID != nil && *r.entities[keys[i]].ParentID == commentID {
			comments = append(comments, r.entities[keys[i]])
		}
		i++
	}

	return comments, nil
}

// GetByPostID returns all comments for a post.
func (r *CommentInMemoryRepository) GetByPostID(ctx context.Context, postID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var keys []uuid.UUID
	for key := range r.entities {
		keys = append(keys, key)
	}

	var comments []*domain.Comment
	i := offset
	for i < len(keys) && len(comments) < limit {
		if r.entities[keys[i]].PostID == postID && r.entities[keys[i]].ParentID == nil {
			comments = append(comments, r.entities[keys[i]])
		}
		i++
	}
	return comments, nil
}

// GetLastComment returns the last comments of a post.
func (r *CommentInMemoryRepository) GetLastComment(ctx context.Context, postID uuid.UUID, lastSeen time.Time, limit int) ([]*domain.Comment, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var keys []uuid.UUID
	for key := range r.entities {
		keys = append(keys, key)
	}

	var comments []*domain.Comment
	i := len(keys) - 1
	for i >= 0 && len(comments) < limit {
		if r.entities[keys[i]].PostID == postID && r.entities[keys[i]].CreatedAt.After(lastSeen) {
			comments = append(comments, r.entities[keys[i]])
		}
		i--
	}
	return comments, nil
}
