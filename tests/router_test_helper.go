package tests

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

// setupTestRouterWithHandlers creates a router with all handlers
// This is a workaround since we can't directly import from main package
// For actual testing, you may need to either:
// 1. Export setupREST from cmd/server or create a test helper there
// 2. Run integration tests against a running server
// 3. Use a build tag to make handlers testable
func setupTestRouterWithHandlers(db *gorm.DB) http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Note: To actually test routes, you need to call setupREST from cmd/server
	// This requires either:
	// - Moving setupREST to an internal package
	// - Using build tags
	// - Running integration tests
	// For now, this is a placeholder that shows the structure

	return router
}
