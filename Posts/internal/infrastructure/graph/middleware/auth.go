package middleware

import (
	"Posts/pkg/jwtservice"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

const userIDKey key = "userID"

// Auth is a middleware that checks if the user is authenticated.
func Auth(jwtGen *jwtservice.Service, logger *slog.Logger) func(next http.Handler) http.Handler {
	const op = "Auth"

	logger = logger.With(slog.Any("op", op))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get the user ID from the request
			token := r.Header.Get("Authorization")

			if token == "" {
				// If there is no token, return an error
				http.Error(w, "no token provided", http.StatusUnauthorized)
				return
			}

			// Check if the token is in the correct format
			if !strings.HasPrefix(token, "Bearer ") {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			// Remove the "Bearer " prefix from the token
			token = strings.TrimPrefix(token, "Bearer ")

			// Get the user ID from the token
			parsedToken, err := jwtGen.Parse(token)
			if err != nil {
				logger.Error("Failed to parse token", slog.String("error", err.Error()))
				http.Error(w, "failed to parse token", http.StatusUnauthorized)
				return
			}

			userID, err := parsedToken.Claims.GetSubject()
			if err != nil {
				logger.Error("Failed to get user ID from token", slog.String("error", err.Error()))
				http.Error(w, "failed to get user ID from token", http.StatusUnauthorized)
				return
			}

			// Add the user ID to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, userIDKey, userID)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID returns the user ID from the request context.
func GetUserID(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}
