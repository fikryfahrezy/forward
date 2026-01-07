# Blog API

Live URL: http://forward-blog-api.fahrezy.work

## Setup and running instructions

### Prerequisites

- Go 1.25 or higher
- Make (optional, for build automation)

#### Local Development

1. **Clone the repository**:

```bash
git clone <repository-url>
cd <clone directory>
```

2. **Install dependencies**:

```bash
go mod download
```

3. **Setup environment**:

```bash
cp .env.example .env
# Edit values on the .env with the real value if needed
```

4. **Generate API documentation**:

```bash
make swagger
```

5. **Start the server**:

```bash
make run
```

### Testing

```bash
# Run all tests
make test
```

### Build the Application

#### Without Docker

```bash
go build -o <executable name> ./cmd/api
```

#### With Docker

```bash
docker -t <container name:tag> build .
```

#### Run with Docker Compose

```bash
docker compose up --build
```

## Architecture explanation

### System Architecture

TBA

### Project Structure

```
project-root/
├── cmd/
│   ├── api/                      # HTTP server entry point with Swagger annotations
│   └── migrate/                  # Database migration CLI for managing database migration
├── internal/                     # Shared packages
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connection management
│   ├── error/                    # Custom error for the project
│   ├── health/                   # Application health status checker
│   ├── logger/                   # Configuration structured logging utilities
│   ├── server/                   # Generic HTTP server with Swagger documentation
│   └── <feature_name>/           # Vertical slicing feature-based modules
│       ├── entity.go             # DTO object and domain model
│       ├── error.go              # Custom error for specific feature
│       ├── handler/              # Presentation layer
│       │   ├── http.go           ## HTTP Route registration
│       │   ├── http_test.go      ## HTTP integration test setup
│       │   ├── <method>.go       ## Handler for the entrypoint
│       │   └── <method>_test.go  ## Handler integration test
│       ├── repository/           # Data access layer
│       │   └── <method>.go       ## Fetching or storing data to persitance or external resource
│       └── service/              # Business logic layer
│           └── <method>.go       ## Main logic for the usecase, handling request from presentation layer and integrating repository
├── docs/                         # Generated Swagger documentation
└── Makefile                      # Build and development commands
```

## Technology choices justification

- Go, ...
- slog, ...
- net/http, ...
- PostgreSQL, ...
- Swagger, ...
- ory/dockertest, ...
- Docker, ...

## Known limitations and future improvements

TBA

## Test coverage report

TBA
