package sql

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/repository/sql/entities"
	"Posts/pkg/logger/slogdiscard"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func setupCommentSQLRepository(t *testing.T) (*gorm.DB, *CommentSQLRepository) {
	slogger := slogdiscard.NewDiscardLogger()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&entities.Comment{})
	if err != nil {
		return nil, nil
	}

	commentRepo := NewCommentSQLRepository(db, slogger)

	return db, commentRepo
}

func TestCommentSQLRepository_GetByPostID(t *testing.T) {
	_, commentRepo := setupCommentSQLRepository(t)

	postID := uuid.New()
	authorID := uuid.New()
	commentID := uuid.New()

	err := commentRepo.Create(context.Background(), &domain.Comment{
		ID:       commentID,
		PostID:   postID,
		AuthorID: authorID,
		Content:  "Test comment",
	})

	if err != nil {
		t.Fatal(err)
	}

	comments, err := commentRepo.GetByPostID(context.Background(), postID, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Equal(t, 1, len(comments))
	assert.Equal(t, commentID, comments[0].ID)
}

func TestCommentSQLRepository_GetChildren(t *testing.T) {
	_, commentRepo := setupCommentSQLRepository(t)

	postID := uuid.New()
	authorID := uuid.New()
	commentID := uuid.New()
	childCommentID := uuid.New()

	err := commentRepo.Create(context.Background(), &domain.Comment{
		ID:       commentID,
		PostID:   postID,
		AuthorID: authorID,
		Content:  "Test comment",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = commentRepo.Create(context.Background(), &domain.Comment{
		ID:       childCommentID,
		PostID:   postID,
		AuthorID: authorID,
		Content:  "Test child comment",
		ParentID: &commentID,
	})
	if err != nil {
		t.Fatal(err)
	}

	comments, err := commentRepo.GetChildren(context.Background(), commentID, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Equal(t, 1, len(comments))
	assert.Equal(t, childCommentID, comments[0].ID)
}

func TestCommentSQLRepository_GetLastComment(t *testing.T) {
	_, commentRepo := setupCommentSQLRepository(t)

	postID := uuid.New()
	authorID := uuid.New()
	commentID := uuid.New()
	lastSeen := time.Now().Add(-time.Hour)
	commentCreatedAt := time.Now()

	err := commentRepo.Create(context.Background(), &domain.Comment{
		ID:        commentID,
		PostID:    postID,
		AuthorID:  authorID,
		Content:   "Test comment",
		CreatedAt: commentCreatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	comments, err := commentRepo.GetLastComment(context.Background(), postID, lastSeen, 10)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Equal(t, 1, len(comments))
	assert.Equal(t, commentID, comments[0].ID)
}
