package inmemory

import (
	"Posts/internal/domain"
	"Posts/internal/interfaces/usecases"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"sync"
)

var _ usecases.AbstractRepositoryInterface[domain.Model] = &AbstractInMemoryRepository[domain.Model]{}

// AbstractInMemoryRepository is a repository for in-memory databases.
type AbstractInMemoryRepository[T domain.Model] struct {
	entities map[uuid.UUID]T
	m        sync.RWMutex
	logger   *slog.Logger
}

// NewAbstractInMemoryRepository creates a new AbstractInMemoryRepository.
func NewAbstractInMemoryRepository[T domain.Model](logger *slog.Logger) AbstractInMemoryRepository[T] {
	return AbstractInMemoryRepository[T]{
		entities: make(map[uuid.UUID]T),
		m:        sync.RWMutex{},
		logger:   logger,
	}
}

// Create creates a new entity.
func (r *AbstractInMemoryRepository[T]) Create(ctx context.Context, entity T) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.entities[entity.GetID()]; ok {
		return domain.ErrAlreadyExists
	}

	r.entities[entity.GetID()] = entity
	return nil
}

// Update updates an entity.
func (r *AbstractInMemoryRepository[T]) Update(ctx context.Context, entity T) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.entities[entity.GetID()]; !ok {
		return domain.ErrNotFound
	}
	r.entities[entity.GetID()] = entity
	return nil
}

// Delete deletes an entity.
func (r *AbstractInMemoryRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.entities[id]; !ok {
		return nil
	}
	delete(r.entities, id)
	return nil
}

// GetByID returns an entity by ID.
func (r *AbstractInMemoryRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (T, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if entity, ok := r.entities[id]; ok {
		return entity, nil
	}
	var entity T
	return entity, domain.ErrNotFound
}

// GetByIds returns entities by IDs.
func (r *AbstractInMemoryRepository[T]) GetByIds(ctx context.Context, ids []uuid.UUID) ([]T, error) {
	const op = "AbstractInMemoryRepository.GetByIds"
	r.m.RLock()
	defer r.m.RUnlock()
	entities := make([]T, 0, len(ids))
	for _, id := range ids {
		if entity, ok := r.entities[id]; ok {
			entities = append(entities, entity)
		}
	}
	r.logger.Debug(op, slog.Any("ids", ids), slog.Any("entities", entities))
	return entities, nil
}

// GetAll returns all entities.
func (r *AbstractInMemoryRepository[T]) GetAll(ctx context.Context, limit int, offset int) ([]T, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	keys := make([]uuid.UUID, 0)
	for key := range r.entities {
		keys = append(keys, key)
	}

	i := offset
	entities := make([]T, 0)
	for i < len(keys) && len(entities) < limit {
		entities = append(entities, r.entities[keys[i]])
		i++
	}
	return entities, nil
}
