package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

// RequireAuth is a middleware that validates the Authorization header token
// and adds the authenticated user to the request context
func RequireAuth(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "authorization header is required")
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeJSONError(w, http.StatusUnauthorized, "authorization header must be in format: Bearer <token>")
				return
			}

			token := parts[1]
			if token == "" {
				writeJSONError(w, http.StatusUnauthorized, "token is required")
				return
			}

			// Get session from database
			session, err := store.GetSessionByToken(db, token)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Check if session is expired
			if time.Now().After(session.ExpiresAt) {
				writeJSONError(w, http.StatusUnauthorized, "token has expired")
				return
			}

			// Get user from session
			if session.UserID == nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid session")
				return
			}

			user, err := store.GetUserByID(db, uint(*session.UserID))
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "user not found")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	return user, ok
}
