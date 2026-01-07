// @title Simple Blog API
// @description A simple API for CRUD Blog & Comment written in Go
// @version 1.0
// @host localhost:8080
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	commentHandler "github.com/fikryfahrezy/forward/blog-api/internal/comment/handler"
	commentRepo "github.com/fikryfahrezy/forward/blog-api/internal/comment/repository"
	commentService "github.com/fikryfahrezy/forward/blog-api/internal/comment/service"
	"github.com/fikryfahrezy/forward/blog-api/internal/config"
	"github.com/fikryfahrezy/forward/blog-api/internal/database"
	"github.com/fikryfahrezy/forward/blog-api/internal/health"
	"github.com/fikryfahrezy/forward/blog-api/internal/logger"
	postHandler "github.com/fikryfahrezy/forward/blog-api/internal/post/handler"
	postRepo "github.com/fikryfahrezy/forward/blog-api/internal/post/repository"
	postService "github.com/fikryfahrezy/forward/blog-api/internal/post/service"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	userHandler "github.com/fikryfahrezy/forward/blog-api/internal/user/handler"
	userRepo "github.com/fikryfahrezy/forward/blog-api/internal/user/repository"
	userService "github.com/fikryfahrezy/forward/blog-api/internal/user/service"

	_ "github.com/fikryfahrezy/forward/blog-api/docs"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Logger, os.Stdout)

	log.Info("Starting application",
		slog.String("server_host", cfg.Server.Host),
		slog.Int("server_port", cfg.Server.Port),
	)

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Initialize repositories
	userRepository := userRepo.New(db.Pool)
	postRepository := postRepo.New(db.Pool)
	commentRepository := commentRepo.New(db.Pool)

	// Initialize services
	userSvc := userService.New(userRepository)
	postSvc := postService.New(postRepository)
	commentSvc := commentService.New(commentRepository)

	// Initialize handlers
	healthHdl := health.NewHealthHandler(db)
	userHdl := userHandler.New(userSvc, cfg.JWT.SecretKey, cfg.JWT.TokenDuration)
	postHdl := postHandler.New(postSvc)
	commentHdl := commentHandler.New(commentSvc)

	// Initialize server
	srv := server.New(server.Config{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	})

	// Configure JWT middleware
	jwtMiddleware := server.NewJWTMiddleware(server.JWTConfig{
		SecretKey: cfg.JWT.SecretKey,
	})
	srv.SetJWTMiddleware(jwtMiddleware)

	// Register route handlers
	routeHandlers := []server.RouteHandler{
		healthHdl,
		userHdl,
		postHdl,
		commentHdl,
	}

	go func() {
		if err := srv.Start(routeHandlers); err != nil {
			log.Error("Server error",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop server
	if err := srv.Stop(ctx); err != nil {
		log.Error("Failed to shutdown server gracefully",
			slog.String("error", err.Error()),
		)
	}

	// Close database connection
	db.Close()
	log.Info("Application shutdown complete")
}
