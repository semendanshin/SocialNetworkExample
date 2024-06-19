package inmemory

import (
	"Posts/internal/domain"
	"Posts/internal/usecases"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

var _ usecases.PostRepository = &PostInMemoryRepository{}

// PostInMemoryRepository is a repository for posts.
type PostInMemoryRepository struct {
	AbstractInMemoryRepository[*domain.Post]
}

// NewPostInMemoryRepository creates a new PostInMemoryRepository.
func NewPostInMemoryRepository(logger *slog.Logger) *PostInMemoryRepository {
	return &PostInMemoryRepository{
		AbstractInMemoryRepository: NewAbstractInMemoryRepository[*domain.Post](logger),
	}
}

// GetByAuthorID returns all posts by a user.
func (r *PostInMemoryRepository) GetByAuthorID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*domain.Post, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var posts []*domain.Post

	var keys []uuid.UUID
	for key := range r.entities {
		keys = append(keys, key)
	}

	i := offset
	for i < len(keys) && len(posts) < limit {
		if r.entities[keys[i]].AuthorID == userID {
			posts = append(posts, r.entities[keys[i]])
		}
		i++
	}

	return posts, nil
}
