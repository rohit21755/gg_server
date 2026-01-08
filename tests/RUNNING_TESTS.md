# How to Run Tests

## Current Status

⚠️ **Note**: The test files currently contain test function stubs with TODO comments. They will compile and run, but they only log messages and don't perform actual assertions yet.

## Basic Commands

### Run All Tests
```bash
cd /Users/rohitkumar/Desktop/gg_ap/backend
go test ./tests/...
```

### Run Specific Test File
```bash
go test ./tests/auth_test.go
go test ./tests/user_test.go
```

### Run with Verbose Output
```bash
go test -v ./tests/...
```

### Run Specific Test Function
```bash
go test -v ./tests/... -run TestHealthCheck
go test -v ./tests/... -run TestLogin
```

### Run with Coverage
```bash
go test -cover ./tests/...
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

### Run Tests in Parallel
```bash
go test -parallel 4 ./tests/...
```

## Prerequisites

Before running functional tests, you need to:

### 1. Set Up Test Database

Create a test database and configure it in `helpers_test.go` or via environment variables:

```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=test_db
export DB_PORT=5432
export JWT_SECRET=test-secret-key
export JWT_REFRESH=test-refresh-secret
export SERVER_PORT=8080
```

### 2. Run Database Migrations

Make sure your test database has all the required tables:

```bash
# Run migrations on test database
# (Adjust based on your migration tool)
```

### 3. Make Router Testable

The tests need access to the router setup. You have a few options:

#### Option A: Export setupREST (Recommended)

Modify `cmd/server/rest.go` to export `setupREST`:

```go
// Change from:
func setupREST(r chi.Router, db *gorm.DB) {

// To:
func SetupREST(r chi.Router, db *gorm.DB) {
```

Then update `tests/helpers_test.go`:

```go
import (
    "github.com/rohit21755/gg_server.git/cmd/server"
    // ... other imports
)

func setupTestRouterWithHandlers(db *gorm.DB) http.Handler {
    router := chi.NewRouter()
    // ... middleware setup
    server.SetupREST(router, db)  // Now accessible
    return router
}
```

#### Option B: Create Test Helper in cmd/server

Create `cmd/server/test_helper.go`:

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"
)

func SetupTestRouter(db *gorm.DB) http.Handler {
    router := chi.NewRouter()
    // ... middleware setup
    setupREST(router, db)
    return router
}
```

Then import it in tests (requires build tag or moving to internal package).

#### Option C: Integration Tests

Run tests against a running server instance (slower but doesn't require code changes).

## Example: Running a Working Test

Once you've implemented a test (see `example_test.go`), you can run it:

```bash
# Run a specific test
go test -v ./tests/... -run TestHealthCheck

# Expected output (once implemented):
# === RUN   TestHealthCheck
# --- PASS: TestHealthCheck (0.01s)
# PASS
# ok      github.com/rohit21755/gg_server.git/tests    0.123s
```

## Troubleshooting

### "package tests is not in GOROOT or GOPATH"
Make sure you're running from the project root directory.

### "undefined: setupTestRouter"
You need to implement the router setup as described above.

### Database Connection Errors
- Ensure PostgreSQL is running
- Check database credentials in environment variables
- Verify test database exists

### Import Errors
- Run `go mod tidy` to ensure dependencies are up to date
- Check that all imports are correct

## Next Steps

1. **Make router testable** (choose one of the options above)
2. **Implement test cases** - Replace TODO comments with actual test logic
3. **Set up CI/CD** - Add tests to your CI pipeline
4. **Add test data fixtures** - Create seed data for consistent testing

## Quick Start Example

Here's a minimal example to get started:

1. Export `setupREST` as `SetupREST` in `cmd/server/rest.go`
2. Update `tests/helpers_test.go` to use `server.SetupREST`
3. Implement one test (e.g., `TestHealthCheck`) following `example_test.go`
4. Run: `go test -v ./tests/... -run TestHealthCheck`
