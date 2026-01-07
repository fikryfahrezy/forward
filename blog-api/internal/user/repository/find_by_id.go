package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password,
			created_at,
			updated_at
		FROM users
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	u := user.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, nil
	}
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}
