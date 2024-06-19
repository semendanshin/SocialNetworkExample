package mappers

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/graph/model"
	"Posts/internal/infrastructure/repository/sql/entities"
)

// ModelToDomainComment maps a model.Comment to a domain.Comment.
func ModelToDomainComment(dto *model.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        dto.ID,
		PostID:    dto.PostID,
		ParentID:  dto.ParentID,
		Content:   dto.Content,
		AuthorID:  dto.AuthorID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// DomainToModelComment maps a domain.Comment to a model.Comment.
func DomainToModelComment(domain *domain.Comment) *model.Comment {
	return &model.Comment{
		ID:        domain.ID,
		PostID:    domain.PostID,
		ParentID:  domain.ParentID,
		Content:   domain.Content,
		AuthorID:  domain.AuthorID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

// DomainToEntityComment maps a domain.Comment to an entities.Comment.
func DomainToEntityComment(domain *domain.Comment) *entities.Comment {
	return &entities.Comment{
		ID:        domain.ID,
		PostID:    domain.PostID,
		ParentID:  domain.ParentID,
		Content:   domain.Content,
		AuthorID:  domain.AuthorID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

// EntityToDomainComment maps an entities.Comment to a domain.Comment.
func EntityToDomainComment(entity *entities.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        entity.ID,
		PostID:    entity.PostID,
		ParentID:  entity.ParentID,
		Content:   entity.Content,
		AuthorID:  entity.AuthorID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// CreateDTOToDomainComment maps a model.NewComment to a domain.Comment.
func CreateDTOToDomainComment(dto *model.NewComment) *domain.Comment {
	return &domain.Comment{
		PostID:   dto.PostID,
		ParentID: dto.ParentID,
		Content:  dto.Content,
		AuthorID: dto.AuthorID,
	}
}
