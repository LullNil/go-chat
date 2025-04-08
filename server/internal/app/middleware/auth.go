package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"go-chat/internal/lib/jwtutil"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/lib/response"
)

// Type for context key to avoid collisions
type contextKey string

const UserIDKey contextKey = "userID"

// Authenticate verifies the JWT token (from the "token" cookie) and adds
// the userID to the context.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				response.RespondError(w, http.StatusUnauthorized,
					"Authentication required")
				return
			}
			response.RespondError(w, http.StatusBadRequest,
				"Invalid authentication cookie")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			response.RespondError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Verify token
		userID, err := jwtutil.VerifyToken(r.Context(), tokenString)
		if err != nil {
			slog.Error("token verification failed", sl.Err(err))
			response.RespondError(w, http.StatusUnauthorized,
				"Invalid or expired token")
			return
		}

		// If the token is valid, add the userID to the context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		// Pass control to the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext helper for extracting userID from context in handlers.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
