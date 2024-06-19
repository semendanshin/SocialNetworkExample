package usecases

import (
	"Posts/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=UserUseCase

// UserUseCase is a use case for users.
type UserUseCase interface {
	AbstractUseCaseInterface[*domain.User]
}
