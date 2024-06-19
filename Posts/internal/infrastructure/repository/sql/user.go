package sql

import (
	"Posts/internal/domain"
	"Posts/internal/infrastructure/repository/sql/entities"
	"Posts/internal/usecases"
	"Posts/internal/utils/mappers"
	"gorm.io/gorm"
	"log/slog"
)

var _ usecases.UserRepository = &UserSQLRepository{}

// UserSQLRepository is a repository for users.
type UserSQLRepository struct {
	AbstractSQLRepository[*domain.User, entities.User]
}

// NewUserSQLRepository creates a new UserSQLRepository.
func NewUserSQLRepository(db *gorm.DB, logger *slog.Logger) *UserSQLRepository {
	return &UserSQLRepository{
		AbstractSQLRepository: NewAbstractSQLRepository[*domain.User, entities.User](
			db, logger, mappers.DomainToEntityUser, mappers.EntityToDomainUser,
		),
	}
}
