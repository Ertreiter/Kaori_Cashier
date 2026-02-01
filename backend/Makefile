.PHONY: dev build run test clean docker-build docker-up docker-down

# Development
dev:
	go run cmd/server/main.go

# Build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server

# Run built binary
run:
	./bin/server

# Test
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

# Database migrations (for local development)
migrate-up:
	@echo "Migrations are run automatically on startup"

# Generate password hash (for creating users)
hash-password:
	@read -p "Enter password: " pwd; \
	go run -e 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { h, _ := bcrypt.GenerateFromPassword([]byte("'$$pwd'"), 10); fmt.Println(string(h)) }'

# Lint
lint:
	golangci-lint run

# Format
fmt:
	go fmt ./...
