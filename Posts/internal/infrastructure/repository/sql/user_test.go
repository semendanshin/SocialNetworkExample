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
)

func setupUserSQLRepository(t *testing.T) (*gorm.DB, *UserSQLRepository) {
	slogger := slogdiscard.NewDiscardLogger()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		return nil, nil
	}

	userRepo := NewUserSQLRepository(db, slogger)

	return db, userRepo
}

func TestUserSQLRepository_Create_Success(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	if err != nil {
		t.Fatal(err)
	}

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test user", user.Name)
	assert.Equal(t, ID, user.ID)
}

func TestUserSQLRepository_Create_AlreadyExists(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestUserSQLRepository_Update_Success(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = userRepo.Update(context.Background(), &domain.User{
		ID:   ID,
		Name: "Updated user",
	})
	if err != nil {
		t.Fatal(err)
	}

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Updated user", user.Name)
	assert.Equal(t, ID, user.ID)
}

func TestUserSQLRepository_Delete_Success(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = userRepo.Delete(context.Background(), ID)
	if err != nil {
		t.Fatal(err)
	}

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, user)
}

func TestUserSQLRepository_Delete_NotFound(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Delete(context.Background(), ID)

	assert.NoError(t, err)
}

func TestUserSQLRepository_GetByID_NotFound(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	user, err := userRepo.GetByID(context.Background(), ID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, user)
}

func TestUserSQLRepository_GetByID_Success(t *testing.T) {
	_, userRepo := setupUserSQLRepository(t)

	ID := uuid.New()
	err := userRepo.Create(context.Background(), &domain.User{
		ID:   ID,
		Name: "Test user",
	})
	if err != nil {
		t.Fatal(err)
	}

	user, err := userRepo.GetByID(context.Background(), ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test user", user.Name)
	assert.Equal(t, ID, user.ID)
}
