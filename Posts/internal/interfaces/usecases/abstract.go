package usecases

import (
	"Posts/internal/domain"
	"context"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=AbstractRepositoryInterface

// AbstractRepositoryInterface is an interface for repositories.
type AbstractRepositoryInterface[T domain.Model] interface {
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (T, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]T, error)
	GetAll(ctx context.Context, limit int, offset int) ([]T, error)
}

// AbstractUseCaseInterface is an interface for use cases.
type AbstractUseCaseInterface[T domain.Model] interface {
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (T, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]T, error)
	GetAll(ctx context.Context, limit int, offset int) ([]T, error)
}

var _ AbstractUseCaseInterface[domain.Model] = &AbstractUseCase[domain.Model]{}

// AbstractUseCase is an abstract use case.
type AbstractUseCase[T domain.Model] struct {
	repository AbstractRepositoryInterface[T]
}

// NewAbstractUseCase creates a new AbstractUseCase.
func NewAbstractUseCase[T domain.Model](repository AbstractRepositoryInterface[T]) AbstractUseCase[T] {
	return AbstractUseCase[T]{repository: repository}
}

// Create creates a new entity.
func (uc *AbstractUseCase[T]) Create(ctx context.Context, entity T) error {
	entity.SetID(uuid.New())
	return uc.repository.Create(ctx, entity)
}

// Update updates an entity.
func (uc *AbstractUseCase[T]) Update(ctx context.Context, entity T) error {
	return uc.repository.Update(ctx, entity)
}

// Delete deletes an entity.
func (uc *AbstractUseCase[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repository.Delete(ctx, id)
}

// GetByID returns an entity by ID.
func (uc *AbstractUseCase[T]) GetByID(ctx context.Context, id uuid.UUID) (T, error) {
	return uc.repository.GetByID(ctx, id)
}

// GetByIds returns entities by IDs.
func (uc *AbstractUseCase[T]) GetByIds(ctx context.Context, ids []uuid.UUID) ([]T, error) {
	return uc.repository.GetByIds(ctx, ids)
}

// GetAll returns all entities.
func (uc *AbstractUseCase[T]) GetAll(ctx context.Context, limit int, offset int) ([]T, error) {
	return uc.repository.GetAll(ctx, limit, offset)
}
