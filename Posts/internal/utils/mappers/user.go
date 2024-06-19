package mappers

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/graph/model"
	"Posts/internal/infrastructure/repository/sql/entities"
)

// ModelToDomainUser maps a model.User to a domain.User.
func ModelToDomainUser(dto *model.User) *domain.User {
	return &domain.User{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

// DomainToModelUser maps a domain.User to a model.User.
func DomainToModelUser(domain *domain.User) *model.User {
	return &model.User{
		ID:   domain.ID,
		Name: domain.Name,
	}
}

// DomainToEntityUser maps a domain.User to an entities.User.
func DomainToEntityUser(domain *domain.User) *entities.User {
	return &entities.User{
		ID:   domain.ID,
		Name: domain.Name,
	}
}

// EntityToDomainUser maps an entities.User to a domain.User.
func EntityToDomainUser(entity *entities.User) *domain.User {
	return &domain.User{
		ID:   entity.ID,
		Name: entity.Name,
	}
}

// CreateDTOToDomainUser maps a model.NewUser to a domain.User.
func CreateDTOToDomainUser(dto *model.NewUser) *domain.User {
	return &domain.User{
		Name: dto.Name,
	}
}
