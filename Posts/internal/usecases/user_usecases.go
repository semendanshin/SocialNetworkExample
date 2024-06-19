package usecases

import (
	"Posts/internal/domain"
	usecaseInterfaces "Posts/internal/interfaces/usecases"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=UserRepository

// UserRepository is a repository for users.
type UserRepository interface {
	usecaseInterfaces.AbstractRepositoryInterface[*domain.User]
}

var _ usecaseInterfaces.UserUseCase = &UserUseCase{}

// UserUseCase is a use case for users.
type UserUseCase struct {
	Repository UserRepository
	usecaseInterfaces.AbstractUseCase[*domain.User]
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(repository UserRepository) *UserUseCase {
	return &UserUseCase{
		Repository:      repository,
		AbstractUseCase: usecaseInterfaces.NewAbstractUseCase[*domain.User](repository),
	}
}
