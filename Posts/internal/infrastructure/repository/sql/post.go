package sql

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/repository/sql/entities"
	"Posts/internal/usecases"
	"Posts/internal/utils/mappers"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
)

var _ usecases.PostRepository = &PostSQLRepository{}

// PostSQLRepository is a repository for posts.
type PostSQLRepository struct {
	AbstractSQLRepository[*domain.Post, entities.Post]
}

// NewPostSQLRepository creates a new PostSQLRepository.
func NewPostSQLRepository(db *gorm.DB, logger *slog.Logger) *PostSQLRepository {
	return &PostSQLRepository{
		AbstractSQLRepository: NewAbstractSQLRepository[*domain.Post, entities.Post](
			db, logger, mappers.DomainToEntityPost, mappers.EntityToDomainPost,
		),
	}
}

// GetByAuthorID returns all posts by a user.
func (r *PostSQLRepository) GetByAuthorID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*domain.Post, error) {
	const op = "PostSQLRepository.GetByAuthorID"
	var posts []*domain.Post
	var postEntities []*entities.Post
	if err := r.db.WithContext(ctx).Where("author_id = ?", userID).Limit(limit).Offset(offset).Find(&postEntities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		return nil, err
	}
	for _, entity := range postEntities {
		posts = append(posts, r.entityToModel(entity))
	}
	return posts, nil

}
