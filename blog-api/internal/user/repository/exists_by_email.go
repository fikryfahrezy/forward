package repository

import (
	"context"
)

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
