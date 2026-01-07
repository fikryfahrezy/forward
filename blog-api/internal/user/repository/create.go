package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (r *Repository) Create(ctx context.Context, u user.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
	)
	return err
}
