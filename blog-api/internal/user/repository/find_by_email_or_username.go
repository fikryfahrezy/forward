package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (r *Repository) FindByEmailOrUsername(ctx context.Context, email, username string) (user.User, error) {
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
			(email = $1 OR username = $2)
			AND deleted_at IS NULL
		LIMIT 1
	`
	u := user.User{}
	err := r.db.QueryRow(ctx, query, email, username).Scan(
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
