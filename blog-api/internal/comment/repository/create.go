package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (r *Repository) Create(ctx context.Context, c *comment.Comment) error {
	query := `
		INSERT INTO comments (
			id,
			content,
			post_id,
			author_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		c.ID,
		c.Content,
		c.PostID,
		c.AuthorID,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return err
}
