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
