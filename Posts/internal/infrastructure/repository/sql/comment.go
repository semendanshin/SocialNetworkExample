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
	"time"
)

var _ usecases.CommentRepository = &CommentSQLRepository{}

// CommentSQLRepository is a repository for comments.
type CommentSQLRepository struct {
	AbstractSQLRepository[*domain.Comment, entities.Comment]
}

// NewCommentSQLRepository creates a new CommentSQLRepository.
func NewCommentSQLRepository(db *gorm.DB, logger *slog.Logger) *CommentSQLRepository {
	return &CommentSQLRepository{
		AbstractSQLRepository: NewAbstractSQLRepository[*domain.Comment, entities.Comment](
			db, logger, mappers.DomainToEntityComment, mappers.EntityToDomainComment,
		),
	}
}

// GetByPostID returns all comments for a post.
func (r *CommentSQLRepository) GetByPostID(ctx context.Context, postID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	const op = "CommentSQLRepository.GetByPostID"
	var comments []*domain.Comment
	var commentEntities []*entities.Comment
	if err := r.db.WithContext(ctx).Where("post_id = ? AND parent_id IS NULL", postID).Limit(limit).Offset(offset).Find(&commentEntities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		return nil, err
	}
	for _, entity := range commentEntities {
		comments = append(comments, r.entityToModel(entity))
	}
	return comments, nil
}

// GetChildren returns all children of a comment.
func (r *CommentSQLRepository) GetChildren(ctx context.Context, commentID uuid.UUID, limit int, offset int) ([]*domain.Comment, error) {
	const op = "CommentSQLRepository.GetChildren"
	var comments []*domain.Comment
	var commentEntities []*entities.Comment
	if err := r.db.WithContext(ctx).Where("parent_id = ?", commentID).Limit(limit).Offset(offset).Find(&commentEntities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		return nil, err
	}
	for _, entity := range commentEntities {
		comments = append(comments, r.entityToModel(entity))
	}
	return comments, nil
}

// GetLastComment returns the last comments of a post.
func (r *CommentSQLRepository) GetLastComment(ctx context.Context, postID uuid.UUID, lastSeen time.Time, limit int) ([]*domain.Comment, error) {
	const op = "CommentSQLRepository.GetLastComment"
	var comments []*domain.Comment
	var commentEntities []*entities.Comment
	if err := r.db.WithContext(ctx).Where("post_id = ? AND created_at > ?", postID, lastSeen).Limit(limit).Find(&commentEntities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		return nil, err
	}
	for _, entity := range commentEntities {
		comments = append(comments, r.entityToModel(entity))
	}
	return comments, nil
}
