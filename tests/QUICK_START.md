# Quick Start: Running Tests

## ✅ Tests Are Ready to Run!

The test files compile and run successfully. Currently, they contain test stubs that log endpoint information.

## Basic Commands

### Run All Tests
```bash
cd /Users/rohitkumar/Desktop/gg_ap/backend
go test ./tests/...
```

### Run with Verbose Output (See What's Being Tested)
```bash
go test -v ./tests/...
```

### Run Specific Test File
```bash
go test -v ./tests/auth_test.go
go test -v ./tests/user_test.go
```

### Run Specific Test Function
```bash
go test -v ./tests/... -run TestHealthCheck
go test -v ./tests/... -run TestLogin
```

### Run with Coverage
```bash
go test -cover ./tests/...
```

### Using Makefile (Easier)
```bash
# Run all tests
make test

# Run with coverage and open HTML report
make test-coverage

# Run specific test
make test-specific TEST=TestHealthCheck
```

## Current Output

When you run the tests now, you'll see output like:
```
=== RUN   TestHealthCheck
    auth_test.go:8: Health check endpoint: GET /api/v1/health
--- PASS: TestHealthCheck (0.00s)
```

This confirms:
- ✅ Tests compile successfully
- ✅ All test functions are being executed
- ✅ Test structure is correct

## Next Steps to Make Tests Functional

1. **Make router testable** (see `RUNNING_TESTS.md` for details)
2. **Implement test logic** - Replace `t.Log()` calls with actual test code
3. **Set up test database** - Configure test DB connection

## Example: Making One Test Work

1. Export `setupREST` in `cmd/server/rest.go`:
   ```go
   func SetupREST(r chi.Router, db *gorm.DB) {  // Changed from setupREST
   ```

2. Update `tests/helpers_test.go` to use it:
   ```go
   import "github.com/rohit21755/gg_server.git/cmd/server"
   // ...
   server.SetupREST(router, db)
   ```

3. Implement `TestHealthCheck` in `tests/auth_test.go`:
   ```go
   func TestHealthCheck(t *testing.T) {
       db := setupTestDB(t)
       router := setupTestRouterWithHandlers(db)
       rr := makeRequest(t, router, "GET", "/api/v1/health", nil, "")
       assertStatusCode(t, rr, http.StatusOK)
   }
   ```

4. Run it:
   ```bash
   go test -v ./tests/... -run TestHealthCheck
   ```

## See Also

- `RUNNING_TESTS.md` - Detailed instructions
- `README.md` - Test structure overview
- `example_test.go` - Reference implementations
