package usecases

import (
	"SSO/internal/domain"
	"context"
	"github.com/google/uuid"
)

type CreateUserDTO struct {
	Username string
	Email    string
	Password string
}

type UpdateUserDTO struct {
	Username string
	Email    string
	Password string
}

type UserUseCases interface {
	Create(ctx context.Context, dto *CreateUserDTO) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, dto *UpdateUserDTO) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, limit int, offset int) ([]*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Login(ctx context.Context, email string, password string) (*domain.User, error)
}
