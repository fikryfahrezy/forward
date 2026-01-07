.PHONY: run test test-clean-cache test-coverage test-coverage-html fmt migrate-up migrate-down migrate-force migrate-version migrate-create swagger


run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-clean-cache:
	go clean -testcache

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

# Migration operations using custom migrate command
migrate-up:
	go run cmd/migrate/main.go -action=up

migrate-down:
	go run cmd/migrate/main.go -action=down

migrate-force:
	go run cmd/migrate/main.go -action=force -version=$(VERSION)

migrate-version:
	go run cmd/migrate/main.go -action=version

migrate-create:
	go run cmd/migrate/main.go -action=create -name=$(NAME)

swagger:
	go tool swag init -g cmd/api/main.go -o docs

# Help
help:
	@echo "Available commands:"
	@echo "  run                - Run the HTTP server"
	@echo "  test               - Run all tests"
	@echo "  test-clean-cache   - Clear test cache"
	@echo "  test-coverage      - Run tests with coverage and show summary"
	@echo "  test-coverage-html - Run tests with coverage and generate HTML report"
	@echo "  fmt                - Format code"
	@echo "  migrate-up         - Run all migrations up"
	@echo "  migrate-down       - Run all migrations down"
	@echo "  migrate-force      - Force migration version (use VERSION=N)"
	@echo "  migrate-version    - Show current migration version"
	@echo "  migrate-create     - Create new migration (use NAME=migration_name)"
	@echo "  swagger            - Generate Swagger documentation"