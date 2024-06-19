package services

import (
	pb "SSO/gen/go"
	"SSO/internal/contracts/usecases"
	"SSO/internal/infrastructure/grpcserver/interceptors"
	"SSO/internal/utils/mappers"
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

var _ pb.UserServiceServer = UserServiceServer{}

// UserServiceServer is used to implement pb.UserServiceServer.
type UserServiceServer struct {
	logger *slog.Logger
	uuc    usecases.UserUseCases
	pb.UnimplementedUserServiceServer
}

// NewUserServiceServer creates a new UserServiceServer.
func NewUserServiceServer(logger *slog.Logger, uuc usecases.UserUseCases) UserServiceServer {
	return UserServiceServer{
		logger: logger,
		uuc:    uuc,
	}
}

// GetMe gets the current user.
func (u UserServiceServer) GetMe(ctx context.Context, empty *emptypb.Empty) (*pb.User, error) {
	userIDStr := interceptors.GetUserID(ctx)
	if userIDStr == "" {
		return nil, errors.New("user id not found")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	user, err := u.uuc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mappers.UserDomainToUserResponse(user), nil
}

// GetUser gets a user by ID.
func (u UserServiceServer) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	userID, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	user, err := u.uuc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mappers.UserDomainToUserResponse(user), nil
}

// ListUsers lists users.
func (u UserServiceServer) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := u.uuc.GetAll(ctx, int(request.Limit), int(request.Offset))
	if err != nil {
		return nil, err

	}

	var usersResponse []*pb.User
	for _, user := range users {
		usersResponse = append(usersResponse, mappers.UserDomainToUserResponse(user))

	}

	return &pb.ListUsersResponse{Users: usersResponse}, nil
}

// CreateUser creates a user.
func (u UserServiceServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.User, error) {
	userDto := mappers.CreateUserRequestToUserDTO(request)
	user, err := u.uuc.Create(ctx, userDto)
	if err != nil {
		return nil, err
	}

	return mappers.UserDomainToUserResponse(user), nil
}

// UpdateUser updates a user.
func (u UserServiceServer) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	userIDStr := interceptors.GetUserID(ctx)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	userDto := mappers.UpdateUserRequestToUserDTO(request)
	user, err := u.uuc.Update(ctx, userID, userDto)
	if err != nil {
		return nil, err
	}

	return mappers.UserDomainToUserResponse(user), nil
}

// DeleteUser deletes a user.
func (u UserServiceServer) DeleteUser(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	userIDStr := interceptors.GetUserID(ctx)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	err = u.uuc.Delete(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
