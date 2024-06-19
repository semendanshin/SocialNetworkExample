package resolvers

import (
	usecaseInterfaces "Posts/internal/interfaces/usecases"
	"log/slog"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is used as dependency injection
type Resolver struct {
	puc    usecaseInterfaces.PostUseCase
	cuc    usecaseInterfaces.CommentUseCase
	uuc    usecaseInterfaces.UserUseCase
	logger *slog.Logger
}

// NewResolver returns a new Resolver
func NewResolver(
	puc usecaseInterfaces.PostUseCase,
	cuc usecaseInterfaces.CommentUseCase,
	uuc usecaseInterfaces.UserUseCase,
	logger *slog.Logger,
) *Resolver {
	return &Resolver{
		puc:    puc,
		cuc:    cuc,
		uuc:    uuc,
		logger: logger,
	}
}
