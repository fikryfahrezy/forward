package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fikryfahrezy/forward/blog-api/internal/config"
	"github.com/fikryfahrezy/forward/blog-api/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var (
		action  = flag.String("action", "", "Migration action: up, down, force, version, create")
		version = flag.Int("version", 0, "Version number for force action")
		name    = flag.String("name", "", "Migration name for create action")
		steps   = flag.Int("steps", 0, "Number of steps for up/down (0 = all)")
	)
	flag.Parse()

	if *action == "" {
		fmt.Println("Usage: migrate -action=<up|down|force|version|create> [options]")
		fmt.Println("Options:")
		fmt.Println("  -version=N    Force to specific version")
		fmt.Println("  -name=NAME    Migration name for create")
		fmt.Println("  -steps=N      Number of steps (0 = all)")
		os.Exit(1)
	}

	cfg := config.Load()
	log := logger.NewLogger(cfg.Logger, os.Stdout)

	if *action == "create" {
		if *name == "" {
			log.Error("Migration name is required for create action")
			os.Exit(1)
		}
		if err := createMigration(*name); err != nil {
			log.Error("Failed to create migration", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		log.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		log.Error("Failed to create database driver", slog.String("error", err.Error()))
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Error("Failed to initialize migrator", slog.String("error", err.Error()))
		os.Exit(1)
	}

	switch *action {
	case "up":
		if err := runUp(m, *steps, log); err != nil {
			log.Error("Migration up failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case "down":
		if err := runDown(m, *steps, log); err != nil {
			log.Error("Migration down failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case "force":
		if *version == 0 {
			log.Error("Version is required for force action")
			os.Exit(1)
		}
		if err := runForce(m, *version, log); err != nil {
			log.Error("Migration force failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case "version":
		if err := showVersion(m, log); err != nil {
			log.Error("Failed to get version", slog.String("error", err.Error()))
			os.Exit(1)
		}
	default:
		log.Error("Invalid action", slog.String("action", *action))
		os.Exit(1)
	}
}

func runUp(m *migrate.Migrate, steps int, log *slog.Logger) error {
	log.Info("Running migrations up", slog.Int("steps", steps))

	var err error
	if steps > 0 {
		err = m.Steps(steps)
	} else {
		err = m.Up()
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Info("No migrations to apply")
	} else {
		log.Info("Migrations applied successfully")
	}

	return nil
}

func runDown(m *migrate.Migrate, steps int, log *slog.Logger) error {
	log.Info("Running migrations down", slog.Int("steps", steps))

	var err error
	if steps > 0 {
		err = m.Steps(-steps)
	} else {
		err = m.Down()
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Info("No migrations to rollback")
	} else {
		log.Info("Migrations rolled back successfully")
	}

	return nil
}

func runForce(m *migrate.Migrate, version int, log *slog.Logger) error {
	log.Info("Forcing migration version", slog.Int("version", version))

	if err := m.Force(version); err != nil {
		return err
	}

	log.Info("Migration version forced successfully")
	return nil
}

func showVersion(m *migrate.Migrate, log *slog.Logger) error {
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Info("No migrations applied yet")
			return nil
		}
		return err
	}

	log.Info("Current migration version",
		slog.Any("version", version),
		slog.Bool("dirty", dirty))

	return nil
}

func createMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name is required")
	}

	timestamp := time.Now().Unix()
	upFile := fmt.Sprintf("migrations/%d_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("migrations/%d_%s.down.sql", timestamp, name)

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Add your UP migration here\n",
		name, time.Now().Format(time.RFC3339))
	downContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Add your DOWN migration here\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upFile, []byte(upContent), 0o644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0o644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Migration files created:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)

	return nil
}
