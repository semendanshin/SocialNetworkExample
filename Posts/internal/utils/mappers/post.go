package mappers

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/graph/model"
	"Posts/internal/infrastructure/repository/sql/entities"
)

// ModelToDomainPost maps a model.Post to a domain.Post.
func ModelToDomainPost(dto *model.Post) *domain.Post {
	return &domain.Post{
		ID:            dto.ID,
		Title:         dto.Title,
		Content:       dto.Content,
		AuthorID:      dto.AuthorID,
		AllowComments: dto.AllowComments,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

// DomainToModelPost maps a domain.Post to a model.Post.
func DomainToModelPost(domain *domain.Post) *model.Post {
	return &model.Post{
		ID:            domain.ID,
		Title:         domain.Title,
		Content:       domain.Content,
		AuthorID:      domain.AuthorID,
		AllowComments: domain.AllowComments,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}
}

// DomainToEntityPost maps a domain.Post to an entities.Post.
func DomainToEntityPost(domain *domain.Post) *entities.Post {
	return &entities.Post{
		ID:            domain.ID,
		Title:         domain.Title,
		Content:       domain.Content,
		AuthorID:      domain.AuthorID,
		AllowComments: domain.AllowComments,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}
}

// EntityToDomainPost maps an entities.Post to a domain.Post.
func EntityToDomainPost(entity *entities.Post) *domain.Post {
	return &domain.Post{
		ID:            entity.ID,
		Title:         entity.Title,
		Content:       entity.Content,
		AuthorID:      entity.AuthorID,
		AllowComments: entity.AllowComments,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// CreateDTOToDomainPost maps a model.NewPost to a domain.Post.
func CreateDTOToDomainPost(dto *model.NewPost) *domain.Post {
	var allowComments bool
	if dto.AllowComments != nil {
		allowComments = *dto.AllowComments
	}
	return &domain.Post{
		Title:         dto.Title,
		Content:       dto.Content,
		AuthorID:      dto.AuthorID,
		AllowComments: allowComments,
	}
}
