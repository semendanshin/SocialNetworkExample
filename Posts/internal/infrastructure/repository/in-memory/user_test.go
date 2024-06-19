package inmemory

import (
	"Posts/internal/domain"
	"Posts/pkg/logger/slogdiscard"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryUserRepository_Create_Success(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})

	assert.NoError(t, err)

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, ID, user.ID)
}

func TestInMemoryUserRepository_Create_AlreadyExists(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	assert.NoError(t, err)

	err = userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestInMemoryUserRepository_Update_Success(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	assert.NoError(t, err)

	err = userRepo.Update(context.Background(), &domain.User{
		ID:   ID,
		Name: "Updated user",
	})
	assert.NoError(t, err)

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Updated user", user.Name)
}

func TestInMemoryUserRepository_Update_NotFound(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Update(context.Background(), &domain.User{
		ID:   ID,
		Name: "Updated user",
	})
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestInMemoryUserRepository_GetByID_Success(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	assert.NoError(t, err)

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test user", user.Name)

}

func TestInMemoryUserRepository_GetByID_NotFound(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	user, err := userRepo.GetByID(context.Background(), ID)

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, user)
}

func TestInMemoryUserRepository_Delete_Success(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	assert.NoError(t, err)

	err = userRepo.Delete(context.Background(), ID)
	assert.NoError(t, err)

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, user)
}

func TestInMemoryUserRepository_Delete_NotFound(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID := uuid.New()
	err := userRepo.Delete(context.Background(), ID)

	assert.NoError(t, err)
}

func TestInMemoryUserRepository_GetAll(t *testing.T) {
	logger := slogdiscard.NewDiscardLogger()
	userRepo := NewUserInMemoryRepository(logger)

	ID1 := uuid.New()
	ID2 := uuid.New()

	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID1,
		Name: "Test user 1",
	})
	assert.NoError(t, err)

	err = userRepo.Create(context.Background(), &domain.User{
		ID:   ID2,
		Name: "Test user 2",
	})
	assert.NoError(t, err)

	users, err := userRepo.GetAll(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 2, len(users))
}
