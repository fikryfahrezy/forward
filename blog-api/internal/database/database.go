package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type Config struct {
	DSN string
}

func NewDB(config Config) (*DB, error) {
	if config.DSN == "" {
		slog.Error("Database DSN is required")
		return nil, fmt.Errorf("database DSN is required")
	}

	pool, err := pgxpool.New(context.Background(), config.DSN)
	if err != nil {
		slog.Error("Failed to open database connection",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db := &DB{Pool: pool}
	if err := db.Health(context.Background()); err != nil {
		slog.Error("Failed to ping database",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Database connection established successfully")

	return db, nil
}

func (db *DB) Close() {
	slog.Info("Closing database connection")
	db.Pool.Close()
}

func (db *DB) Health(ctx context.Context) error {
	return db.Ping(ctx)
}
