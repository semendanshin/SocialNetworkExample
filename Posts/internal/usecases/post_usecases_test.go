package usecases

import (
	"Posts/internal/usecases/mocks"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPostUseCase_GetByAuthorID(t *testing.T) {
	repo := &mocks.PostRepository{}
	uc := NewPostUseCase(repo)

	repo.On("GetByAuthorID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	id := uuid.New()
	_, err := uc.GetByAuthorID(context.Background(), id, 0, 0)

	assert.NoError(t, err)
}
