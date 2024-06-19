package inmemory

import (
	"Posts/internal/domain"
	"Posts/pkg/logger/slogdiscard"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupPostInMemoryRepository(t *testing.T) *PostInMemoryRepository {
	logger := slogdiscard.NewDiscardLogger()

	return NewPostInMemoryRepository(logger)
}

func TestPostInMemoryRepository_Create_Success(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)
}

func TestPostInMemoryRepository_Create_AlreadyExists(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

	err = rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestPostInMemoryRepository_Update(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

	err = rep.Update(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

}

func TestPostInMemoryRepository_Update_NotFound(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	err := rep.Update(context.Background(), &domain.Post{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.Equal(t, domain.ErrNotFound, err)
}

func TestPostInMemoryRepository_GetByID_Success(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

	post, err := rep.GetByID(context.Background(), postID)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, postID, post.ID)
}

func TestPostInMemoryRepository_GetByID_NotFound(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()

	post, err := rep.GetByID(context.Background(), postID)

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, post)
}

func TestPostInMemoryRepository_Delete_Success(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

	err = rep.Delete(context.Background(), postID)

	assert.NoError(t, err)

	post, err := rep.GetByID(context.Background(), postID)

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, post)
}

func TestPostInMemoryRepository_Delete_NotFound(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()

	err := rep.Delete(context.Background(), postID)

	assert.NoError(t, err)
}

func TestPostInMemoryRepository_GetAll(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	postID := uuid.New()
	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: uuid.New(),
		Title:    "Test post",
		Content:  "Test content",
	})

	assert.NoError(t, err)

	posts, err := rep.GetAll(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, postID, posts[0].ID)
}

func TestPostInMemoryRepository_GetByAuthorID(t *testing.T) {
	rep := setupPostInMemoryRepository(t)

	authorID := uuid.New()
	postID := uuid.New()

	err := rep.Create(context.Background(), &domain.Post{
		ID:       postID,
		AuthorID: authorID,
		Title:    "Test post",
		Content:  "Test content",
	})

	if err != nil {
		t.Fatal(err)
	}

	posts, err := rep.GetByAuthorID(context.Background(), authorID, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, postID, posts[0].ID)
}
