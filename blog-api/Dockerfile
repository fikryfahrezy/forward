# Build stage
FROM golang:1.25.5-bookworm AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate swagger
RUN make swagger

# Build the application using Makefile
RUN go build -o api ./cmd/api

# Build the migrator binary using Makefile
RUN go build -o migrate ./cmd/migrate

# Final stage
FROM ubuntu:24.04

# Install ca-certificates and curl for health checks
RUN apt-get update && \
    apt-get install -y ca-certificates curl tzdata && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -m -s /bin/false appuser

# Set working directory
WORKDIR /app

# Copy api binary from builder stage
COPY --from=builder /app/api .

# Copy migrator binary from builder stage
COPY --from=builder /app/migrate .

# Copy migrations from builder stage
COPY --from=builder /app/migrations ./migrations

# Copy swagger folder from builder stage
COPY --from=builder /app/docs ./docs

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set default port
ENV PORT=8080

# Expose port
EXPOSE $PORT

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT}/api/health || exit 1

# Run the application
CMD ["sh", "-c", "./api"]