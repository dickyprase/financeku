BINARY_NAME=financeku
BUILD_DIR=./cmd/api

.PHONY: build run dev test clean migrate seed docker-up docker-down

build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

run: build
	./$(BINARY_NAME)

dev:
	go run $(BUILD_DIR)/main.go

seed:
	go run $(BUILD_DIR)/main.go --seed

migrate:
	go run ./cmd/migrate/main.go

migrate-seed:
	go run ./cmd/migrate/main.go --seed

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet
