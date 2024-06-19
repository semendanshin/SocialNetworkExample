package usecases

import (
	"SSO/internal/domain"
	"context"
	"crypto/sha256"
	"github.com/google/uuid"

	contractUseCases "SSO/internal/contracts/usecases"
)

// UserRepositoryInterface is an interface for user repositories.
type UserRepositoryInterface interface {
	contractUseCases.AbstractRepositoryInterface[*domain.User]
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

var _ contractUseCases.UserUseCases = &UserUseCases{}

// UserUseCases is a use case for users.
type UserUseCases struct {
	Repository UserRepositoryInterface
	contractUseCases.AbstractUseCases[*domain.User]
}

// NewUserUseCases creates a new UserUseCases.
func NewUserUseCases(repository UserRepositoryInterface) *UserUseCases {
	return &UserUseCases{
		Repository:       repository,
		AbstractUseCases: contractUseCases.NewAbstractUseCase[*domain.User](repository),
	}
}

// GetByUsername gets a user by username.
func (u *UserUseCases) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.Repository.GetByUsername(ctx, username)
}

// GetByEmail gets a user by email.
func (u *UserUseCases) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.Repository.GetByEmail(ctx, email)
}

// Login logs in a user.
func (u *UserUseCases) Login(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := u.Repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}
	if user.HashedPassword != hashPassword(password) {
		return nil, domain.ErrInvalidPassword
	}
	return user, nil
}

// Create creates a new user.
func (u *UserUseCases) Create(ctx context.Context, dto *contractUseCases.CreateUserDTO) (*domain.User, error) {
	user := &domain.User{
		UUID:           uuid.New(),
		Username:       dto.Username,
		Email:          dto.Email,
		HashedPassword: hashPassword(dto.Password),
	}
	return user, u.Repository.Create(ctx, user)
}

// Update updates a user.
func (u *UserUseCases) Update(ctx context.Context, ID uuid.UUID, dto *contractUseCases.UpdateUserDTO) (*domain.User, error) {
	user, err := u.Repository.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}
	switch {
	case dto.Username != "":
		user.Username = dto.Username
	case dto.Email != "":
		user.Email = dto.Email
	case dto.Password != "":
		user.HashedPassword = hashPassword(dto.Password)
	}
	return user, u.Repository.Update(ctx, user)
}

func hashPassword(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
