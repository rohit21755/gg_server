# API Route Tests

This directory contains comprehensive test files for all API routes in the Campus Ambassador Platform.

## Structure

The tests are organized by feature/domain:

- `auth_test.go` - Authentication routes (login, register, refresh, etc.)
- `public_test.go` - Public routes (colleges, states, leaderboards)
- `user_test.go` - User profile and management routes
- `task_test.go` - Task-related routes
- `submission_test.go` - Submission routes
- `campaign_test.go` - Campaign routes
- `gamification_test.go` - XP, levels, badges, streaks, spin wheel
- `engagement_test.go` - Flash challenges, trivia, mystery boxes, battles
- `rewards_test.go` - Rewards and redemptions
- `referral_test.go` - Referral system
- `college_state_test.go` - College and state routes
- `campus_wars_test.go` - Campus wars routes
- `survey_test.go` - Survey routes
- `notification_test.go` - Notification routes
- `wallet_test.go` - Wallet and transactions
- `social_test.go` - Social feed and posts
- `dashboard_test.go` - Dashboard routes
- `email_test.go` - Email preferences
- `admin_test.go` - Admin-only routes
- `helpers_test.go` - Test helper utilities
- `router_test_helper.go` - Router setup helper

## Current Status

The test files currently contain test function stubs with TODO comments. To make these tests fully functional, you need to:

1. **Make the router setup testable**: Either:
   - Export `setupREST` from `cmd/server` package, or
   - Create a test helper in `cmd/server` that can be imported, or
   - Use build tags to make handlers testable

2. **Set up test database**: Configure a test database connection in `helpers_test.go`

3. **Implement test cases**: Replace TODO comments with actual test implementations

## Running Tests

Once implemented, run tests with:

```bash
# Run all tests
go test ./tests/...

# Run specific test file
go test ./tests/auth_test.go

# Run with verbose output
go test -v ./tests/...

# Run with coverage
go test -cover ./tests/...
```

## Test Database Setup

Before running tests, ensure you have:

1. A test database configured (set via environment variables or in `helpers_test.go`)
2. Database migrations run on the test database
3. Test data seeded if needed

## Example Test Implementation

See `example_test.go` for a complete example of how to implement a test once the router setup is made testable.

## Notes

- All protected routes require authentication tokens
- Admin routes require admin role
- Some endpoints return 501 Not Implemented (marked in test comments)
- Tests should clean up test data after execution
