package inmemory

import (
	"Posts/internal/domain"
	"Posts/internal/usecases"
	"log/slog"
)

var _ usecases.UserRepository = &UserInMemoryRepository{}

// UserInMemoryRepository is a repository for users.
type UserInMemoryRepository struct {
	AbstractInMemoryRepository[*domain.User]
}

// NewUserInMemoryRepository creates a new UserInMemoryRepository.
func NewUserInMemoryRepository(logger *slog.Logger) *UserInMemoryRepository {
	return &UserInMemoryRepository{
		AbstractInMemoryRepository: NewAbstractInMemoryRepository[*domain.User](logger),
	}
}
