package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM posts WHERE slug = $1 AND deleted_at IS NULL
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, slug).Scan(&exists)
	return exists, err
}

func (r *Repository) ExistsBySlugExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM posts WHERE slug = $1 AND id != $2 AND deleted_at IS NULL
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, slug, excludeID).Scan(&exists)
	return exists, err
}
