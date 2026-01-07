// @title Simple Blog API
// @description A simple API for CRUD Blog & Comment written in Go
// @version 1.0
// @host localhost:8080
// @BasePath /api
// @schemes http https
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fikryfahrezy/forward/blog-api/internal/config"
	"github.com/fikryfahrezy/forward/blog-api/internal/database"
	"github.com/fikryfahrezy/forward/blog-api/internal/health"
	"github.com/fikryfahrezy/forward/blog-api/internal/http"
	"github.com/fikryfahrezy/forward/blog-api/internal/logger"

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

	healthHandler := health.NewHealthHandler(db)

	srv := http.New(http.Config{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	})
	routeHandlers := []http.RouteHandler{
		healthHandler,
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
