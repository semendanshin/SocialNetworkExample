package gormrepository

import (
	"SSO/internal/contracts/usecases"
	"SSO/internal/domain"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
)

// Entity is an interface for entities.
type Entity interface{}

var _ usecases.AbstractRepositoryInterface[domain.Model] = &AbstractGormRepository[domain.Model, Entity]{}

// AbstractGormRepository is a repository for SQL databases.
type AbstractGormRepository[TModel domain.Model, TEntity Entity] struct {
	db            *gorm.DB
	logger        *slog.Logger
	modelToEntity func(model TModel) *TEntity
	entityToModel func(entity *TEntity) TModel
}

// NewAbstractGormRepository creates a new AbstractGormRepository.
func NewAbstractGormRepository[TModel domain.Model, TEntity Entity](
	db *gorm.DB,
	logger *slog.Logger,
	modelToEntity func(model TModel) *TEntity,
	entityToModel func(entity *TEntity) TModel,
) AbstractGormRepository[TModel, TEntity] {
	return AbstractGormRepository[TModel, TEntity]{
		db:            db,
		logger:        logger,
		modelToEntity: modelToEntity,
		entityToModel: entityToModel,
	}
}

// Create creates a new entity.
func (r *AbstractGormRepository[TModel, TEntity]) Create(ctx context.Context, model TModel) error {
	const op = "AbstractGormRepository.Create"
	entity := r.modelToEntity(model)
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

// Update updates an entity.
func (r *AbstractGormRepository[TModel, TEntity]) Update(ctx context.Context, model TModel) error {
	const op = "AbstractGormRepository.Update"
	entity := r.modelToEntity(model)
	if err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

// Delete deletes an entity.
func (r *AbstractGormRepository[TModel, TEntity]) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "AbstractGormRepository.Delete"
	var entity TEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

// GetByID returns an entity by ID.
func (r *AbstractGormRepository[TModel, TEntity]) GetByID(ctx context.Context, id uuid.UUID) (TModel, error) {
	const op = "AbstractGormRepository.GetByID"
	var entity TEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		var model TModel
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model, domain.ErrNotFound
		}
		return model, err
	}
	return r.entityToModel(&entity), nil
}

// GetByIds returns entities by IDs.
func (r *AbstractGormRepository[TModel, TEntity]) GetByIds(ctx context.Context, ids []uuid.UUID) ([]TModel, error) {
	const op = "AbstractGormRepository.GetByIds"
	var entities []*TEntity
	if err := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&entities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	var models []TModel
	for _, entity := range entities {
		models = append(models, r.entityToModel(entity))
	}
	return models, nil
}

// GetAll returns all entities.
func (r *AbstractGormRepository[TModel, TEntity]) GetAll(ctx context.Context, limit int, offset int) ([]TModel, error) {
	const op = "AbstractGormRepository.GetAll"
	var entities []*TEntity
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		r.logger.Error(op, slog.Any("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	var models []TModel
	for _, entity := range entities {
		models = append(models, r.entityToModel(entity))
	}
	return models, nil
}
