package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) PostExists(ctx context.Context, postID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM posts WHERE id = $1 AND deleted_at IS NULL
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, postID).Scan(&exists)
	return exists, err
}
