package interceptors

import (
	"SSO/pkg/jwt"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"
)

// UserIDKey is the key used to store the user ID in the context.
const UserIDKey key = "user_id"

// UnprotectedMethods is a map of unprotected methods.
var UnprotectedMethods = map[string]struct{}{
	"/AuthService/Login":        {},
	"/AuthService/RefreshToken": {},
	"/AuthService/Logout":       {},

	"/UserService/ListUsers":  {},
	"/UserService/CreateUser": {},
	"/UserService/GetUser":    {},
}

// AuthInterceptor is a gRPC interceptor that checks if the request is authenticated.
func AuthInterceptor(w jwt.Parser, logger *slog.Logger) grpc.UnaryServerInterceptor {
	const op = "AuthInterceptor"

	logger = logger.With(slog.Any("op", op))
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		logger.Debug("Calling method", slog.String("method", info.FullMethod))

		if _, ok := UnprotectedMethods[info.FullMethod]; ok {
			logger.Debug("Method is unprotected", slog.String("method", info.FullMethod))
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		tokens, ok := md["authorization"]
		if !ok || len(tokens) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := tokens[0]

		if token == "" {
			// If there is no token, return an error
			return nil, status.Errorf(codes.Unauthenticated, "no token provided")
		}

		// Check if the token is in the correct format
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
		}

		// Remove the "Bearer " prefix from the token
		token = strings.TrimPrefix(token, "Bearer ")

		logger.Debug("Got auth token", slog.Any("token", token))

		//Get the user ID from the token
		parsedToken, err := w.Parse(token)
		if err != nil {
			logger.Error("Failed to parse token", slog.String("error", err.Error()))
			return nil, status.Errorf(codes.Unauthenticated, "failed to parse token")
		}

		userID, err := parsedToken.Claims.GetSubject()
		if err != nil {
			logger.Error("Failed to get user ID from token", slog.String("error", err.Error()))
			return nil, status.Errorf(codes.Unauthenticated, "failed to get user ID from token")
		}

		// Add the user ID to the request context
		ctx = context.WithValue(ctx, UserIDKey, userID)

		return handler(ctx, req)
	}
}

// GetUserID returns the user ID from the context.
func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(UserIDKey).(string)
	return userID
}
