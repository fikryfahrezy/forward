package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (r *Repository) Update(ctx context.Context, c comment.Comment) error {
	query := `
		UPDATE comments SET
			content = $1,
			updated_at = $2
		WHERE 
			id = $3
			AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query,
		c.Content,
		c.UpdatedAt,
		c.ID,
	)
	return err
}
