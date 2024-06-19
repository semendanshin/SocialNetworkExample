package mappers

import (
	pb "SSO/gen/go"
	"SSO/internal/contracts/usecases"
	"SSO/internal/domain"
	"SSO/internal/infrastructure/repositories/gorm/entities"
	"encoding/hex"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserDomainToUserResponse maps a domain user to a user response.
func UserDomainToUserResponse(domainUser *domain.User) *pb.User {
	return &pb.User{
		Id:        domainUser.UUID.String(),
		Username:  domainUser.Username,
		Email:     domainUser.Email,
		CreatedAt: timestamppb.New(domainUser.CreatedAt),
		UpdatedAt: timestamppb.New(domainUser.UpdatedAt),
	}
}

// CreateUserRequestToUserDTO maps a create user request to a create user DTO.
func CreateUserRequestToUserDTO(request *pb.CreateUserRequest) *usecases.CreateUserDTO {
	return &usecases.CreateUserDTO{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

// UserDomainToUserEntity maps a domain user to a user entity.
func UserDomainToUserEntity(domainUser *domain.User) *entities.User {
	return &entities.User{
		ID:             domainUser.UUID,
		Username:       domainUser.Username,
		Email:          domainUser.Email,
		HashedPassword: hex.EncodeToString(domainUser.HashedPassword[:]),
		CreatedAt:      domainUser.CreatedAt,
		UpdatedAt:      domainUser.UpdatedAt,
	}
}

// UserEntityToUserDomain maps a user entity to a domain user.
func UserEntityToUserDomain(entity *entities.User) *domain.User {
	pass, err := hex.DecodeString(entity.HashedPassword)
	if err != nil {
		panic(err)
	}

	return &domain.User{
		UUID:           entity.ID,
		Username:       entity.Username,
		Email:          entity.Email,
		HashedPassword: [32]byte(pass),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// UpdateUserRequestToUserDTO maps an update user request to an update user DTO.
func UpdateUserRequestToUserDTO(request *pb.UpdateUserRequest) *usecases.UpdateUserDTO {
	return &usecases.UpdateUserDTO{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}
