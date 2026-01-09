package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (r *Repository) Update(ctx context.Context, p post.Post) error {
	query := `
		UPDATE posts SET
			title = $1,
			slug = $2,
			content = $3,
			updated_at = NOW()
		WHERE
			id = $4
			AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query,
		p.Title,
		p.Slug,
		p.Content,
		p.ID,
	)
	return err
}
