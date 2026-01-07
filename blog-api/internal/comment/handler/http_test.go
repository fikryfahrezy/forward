package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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

	commentHandler "github.com/fikryfahrezy/forward/blog-api/internal/comment/handler"
	commentRepository "github.com/fikryfahrezy/forward/blog-api/internal/comment/repository"
	commentService "github.com/fikryfahrezy/forward/blog-api/internal/comment/service"
	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	postHandler "github.com/fikryfahrezy/forward/blog-api/internal/post/handler"
	postRepository "github.com/fikryfahrezy/forward/blog-api/internal/post/repository"
	postService "github.com/fikryfahrezy/forward/blog-api/internal/post/service"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user"
	userHandler "github.com/fikryfahrezy/forward/blog-api/internal/user/handler"
	userRepository "github.com/fikryfahrezy/forward/blog-api/internal/user/repository"
	userService "github.com/fikryfahrezy/forward/blog-api/internal/user/service"
)

var (
	testPool           *pgxpool.Pool
	testCommentHandler *commentHandler.Handler
	testPostHandler    *postHandler.Handler
	testUserHandler    *userHandler.Handler
	testServer         *server.Server
)

const (
	testJWTSecret   = "test-secret-key-for-integration-tests"
	testTokenExpiry = 24 * time.Hour
	testDBUser      = "test"
	testDBPassword  = "test"
	testDBName      = "testdb"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pool.MaxWait = 120 * time.Second

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

	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not set expiration: %s", err)
	}

	hostAndPort := getHostPort(resource, "5432/tcp")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser, testDBPassword, hostAndPort, testDBName)

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

	if err := runMigrations(databaseURL); err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	setupTestHandlers()

	code := m.Run()

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

func setupTestHandlers() {
	userRepo := userRepository.New(testPool)
	userSvc := userService.New(userRepo)
	testUserHandler = userHandler.New(userSvc, testJWTSecret, testTokenExpiry)

	postRepo := postRepository.New(testPool)
	postSvc := postService.New(postRepo)
	testPostHandler = postHandler.New(postSvc)

	commentRepo := commentRepository.New(testPool)
	commentSvc := commentService.New(commentRepo)
	testCommentHandler = commentHandler.New(commentSvc)

	testServer = server.New(server.Config{Host: "localhost", Port: 8080})
	testServer.SetJWTMiddleware(server.NewJWTMiddleware(server.JWTConfig{SecretKey: testJWTSecret}))
	testUserHandler.SetupRoutes(testServer)
	testPostHandler.SetupRoutes(testServer)
	testCommentHandler.SetupRoutes(testServer)
}

func cleanup(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "DELETE FROM comments")
	if err != nil {
		t.Fatalf("Failed to cleanup comments: %v", err)
	}
	_, err = testPool.Exec(context.Background(), "DELETE FROM posts")
	if err != nil {
		t.Fatalf("Failed to cleanup posts: %v", err)
	}
	_, err = testPool.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}

func registerAndGetToken(t *testing.T, username, email, password string) string {
	t.Helper()

	reqBody := user.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to register user: %s", rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	result := response.Result.(map[string]any)
	return result["token"].(string)
}

func createPost(t *testing.T, token, title, content string) string {
	t.Helper()

	reqBody := post.CreatePostRequest{
		Title:   title,
		Content: content,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to create post: %s", rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	result := response.Result.(map[string]any)
	return result["id"].(string)
}
