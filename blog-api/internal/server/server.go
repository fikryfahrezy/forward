package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Config holds the server configuration
type Config struct {
	Host string
	Port int
}

// Server manages the application server lifecycle
type Server struct {
	config        Config
	mux           *http.ServeMux
	server        *http.Server
	jwtMiddleware *JWTMiddleware
}

// RouteHandler defines interface for features to register their routes
type RouteHandler interface {
	SetupRoutes(server *Server)
}

// New creates a new server instance
func New(config Config) *Server {
	// Initialize Mux server
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	return &Server{
		config: config,
		mux:    mux,
		server: server,
	}
}

// Start starts the server
func (s *Server) Start(handlers []RouteHandler) error {
	if s.mux == nil {
		return fmt.Errorf("server not initialized, call New() first")
	}

	slog.Info("Initializing server")

	// Swagger documentation
	s.mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Register all the handlers routes
	for _, handler := range handlers {
		handler.SetupRoutes(s)
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	slog.Info("Starting server",
		slog.String("address", addr),
	)
	return s.server.ListenAndServe()
}

// Stop stops the server gracefully
func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping server...")
	return s.server.Shutdown(ctx)
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.Handle(pattern, LoggerMiddleware(RecoverMiddleware(CORSMiddleware((http.HandlerFunc(handler))))))
}

func (s *Server) HandleFuncWithAuth(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if s.jwtMiddleware == nil {
		panic("JWT middleware not configured. Call SetJWTMiddleware first.")
	}
	s.mux.Handle(pattern, LoggerMiddleware(RecoverMiddleware(CORSMiddleware(s.jwtMiddleware.Middleware(http.HandlerFunc(handler))))))
}

func (s *Server) SetJWTMiddleware(jwtMiddleware *JWTMiddleware) {
	s.jwtMiddleware = jwtMiddleware
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("Started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		next.ServeHTTP(w, r)

		slog.Info("Completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				slog.Info("Panic recovered",
					slog.Any("reason", rcv),
					slog.String("stack trace", string(debug.Stack())),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
