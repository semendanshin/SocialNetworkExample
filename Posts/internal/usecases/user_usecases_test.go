package usecases

import (
	"Posts/internal/domain"
	"Posts/internal/usecases/mocks"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserUseCase_Create(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := NewUserUseCase(repo)

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	entity := &domain.User{}
	err := uc.Create(context.Background(), entity)

	assert.NoError(t, err)
}

func TestUserUseCase_Update(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := NewUserUseCase(repo)

	repo.On("Update", mock.Anything, mock.Anything).Return(nil)

	entity := &domain.User{}
	err := uc.Update(context.Background(), entity)

	assert.NoError(t, err)
}

func TestUserUseCase_Delete(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := NewUserUseCase(repo)

	repo.On("Delete", mock.Anything, mock.Anything).Return(nil)

	id := uuid.UUID{}
	err := uc.Delete(context.Background(), id)

	assert.NoError(t, err)
}

func TestUserUseCase_GetByID(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := NewUserUseCase(repo)

	repo.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)

	id := uuid.UUID{}
	_, err := uc.GetByID(context.Background(), id)

	assert.NoError(t, err)
}

func TestUserUseCase_GetByIds(t *testing.T) {
	repo := &mocks.UserRepository{}
	uc := NewUserUseCase(repo)

	repo.On("GetByIds", mock.Anything, mock.Anything).Return(nil, nil)

	var ids []uuid.UUID
	_, err := uc.GetByIds(context.Background(), ids)

	assert.NoError(t, err)
}
