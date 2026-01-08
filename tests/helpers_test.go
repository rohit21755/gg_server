package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rohit21755/gg_server.git/internal/db"
	"github.com/rohit21755/gg_server.git/internal/env"
	"github.com/rohit21755/gg_server.git/internal/store"
	"github.com/rohit21755/gg_server.git/ws"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter http.Handler

// setupTestDB initializes a test database connection
func setupTestDB(t *testing.T) *gorm.DB {
	if testDB != nil {
		return testDB
	}

	// Set test environment variables if not already set
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PASS", "postgres")
		os.Setenv("DB_NAME", "test_db")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("JWT_SECRET", "test-secret-key")
		os.Setenv("JWT_REFRESH", "test-refresh-secret")
		os.Setenv("SERVER_PORT", "8080")
	}

	env.Load()
	database := db.Connect()
	testDB = database
	return database
}

// setupTestRouter creates a test router with all routes
// This mimics the setup in main.go but is accessible from tests
func setupTestRouter(t *testing.T) http.Handler {
	if testRouter != nil {
		return testRouter
	}

	_ = setupTestDB(t) // Initialize test DB
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

	// WebSockets
	hub := ws.NewHub()
	go hub.Run()

	// REST API - we need to import setupREST from cmd/server
	// For now, we'll create a minimal router that can be extended
	// In a real scenario, you might want to export setupREST or create a test helper
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// WebSocket handler placeholder
	})

	// Note: To fully test routes, you may need to either:
	// 1. Export setupREST from cmd/server package
	// 2. Or create a test-specific router setup
	// For now, tests will need to set up routes manually or use integration tests

	testRouter = router
	return router
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *gorm.DB, email, password, role string) *store.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &store.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    "Test",
		LastName:     "User",
		Role:         role,
		IsActive:     true,
		ReferralCode: fmt.Sprintf("REF%d", time.Now().UnixNano()),
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// createTestSession creates a test session for a user
func createTestSession(t *testing.T, db *gorm.DB, userID uint) string {
	token := generateTestToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	uid := int(userID)

	session := &store.UserSession{
		SessionToken: token,
		UserID:       &uid,
		ExpiresAt:    expiresAt,
		LastActive:   time.Now(),
	}

	if err := store.CreateSession(db, session); err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	return token
}

// generateTestToken generates a random token for testing
func generateTestToken() string {
	return fmt.Sprintf("test-token-%d", time.Now().UnixNano())
}

// makeRequest makes an HTTP request to the test server
func makeRequest(t *testing.T, router http.Handler, method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// parseJSONResponse parses a JSON response
func parseJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, target interface{}) {
	if err := json.Unmarshal(rr.Body.Bytes(), target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nResponse body: %s", err, rr.Body.String())
	}
}

// cleanupTestData cleans up test data from the database
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// Clean up in reverse order of dependencies
	db.Exec("DELETE FROM user_sessions")
	db.Exec("DELETE FROM users")
	// Add more cleanup as needed
}

// assertStatusCode checks if the response has the expected status code
func assertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d\nResponse body: %s", expected, rr.Code, rr.Body.String())
	}
}

// assertJSONError checks if the response contains an error message
func assertJSONError(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assertStatusCode(t, rr, expectedStatus)

	var response struct {
		Error string `json:"error"`
	}
	parseJSONResponse(t, rr, &response)

	if response.Error != expectedError && expectedError != "" {
		t.Errorf("Expected error '%s', got '%s'", expectedError, response.Error)
	}
}
