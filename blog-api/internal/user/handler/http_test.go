package handler_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/fikryfahrezy/forward/blog-api/internal/logger"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user/handler"
	"github.com/fikryfahrezy/forward/blog-api/internal/user/repository"
	"github.com/fikryfahrezy/forward/blog-api/internal/user/service"
)

var (
	testPool    *pgxpool.Pool
	testHandler *handler.Handler
	testServer  *server.Server
)

const (
	testJWTSecret   = "test-secret-key-for-integration-tests"
	testTokenExpiry = 24 * time.Hour
	testDBUser      = "test"
	testDBPassword  = "test"
	testDBName      = "testdb"
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	// Set maximum wait time for container
	pool.MaxWait = 120 * time.Second

	// Start PostgreSQL container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17.4-bookworm",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", testDBUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testDBPassword),
			fmt.Sprintf("POSTGRES_DB=%s", testDBName),
			"listen_addresses='*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Set container to expire after 120 seconds (for cleanup in case of failure)
	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not set expiration: %s", err)
	}

	// Get the database host and port
	hostAndPort := getHostPort(resource, "5432/tcp")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser, testDBPassword, hostAndPort, testDBName)

	// Retry connection until database is ready
	if err := pool.Retry(func() error {
		var err error
		testPool, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return testPool.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// Run migrations
	if err := runMigrations(databaseURL); err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	// Setup handler
	setupTestHandler()

	// Run tests
	code := m.Run()

	// Cleanup
	testPool.Close()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	return u.Hostname() + ":" + resource.GetPort(id)
}

// getMigrationsPath returns the absolute path to the migrations directory.
// It finds the project root by searching for go.mod, making it robust
// regardless of where the test file is located.
func getMigrationsPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get current file path")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "migrations")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("Could not find project root (go.mod)")
		}
		dir = parent
	}
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	//nolint:errcheck
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	migrationsPath := getMigrationsPath()
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func setupTestHandler() {
	logger.NewLogger(logger.Config{}, io.Discard)
	repo := repository.New(testPool)
	svc := service.New(server.NewJWTGenerator(testJWTSecret, testTokenExpiry), repo)
	testHandler = handler.New(svc)

	testServer = server.New(server.Config{Host: "localhost", Port: 8080})
	testServer.SetJWTMiddleware(server.NewJWTMiddleware(server.JWTConfig{SecretKey: testJWTSecret}))
	testHandler.SetupRoutes(testServer)
}

func cleanupUsers(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}
