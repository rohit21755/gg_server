run:
	go run cmd/server/*.go

build:
	go build -o bin/server cmd/server/*.go

docker:
	docker-compose up --build

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/yourapp?sslmode=disable" up

seed:
	go run internal/db/seed.go

test:
	go test -v ./tests/...

test-coverage:
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out

test-specific:
	@echo "Usage: make test-specific TEST=TestName"
	@echo "Example: make test-specific TEST=TestHealthCheck"
	go test -v ./tests/... -run $(TEST)
