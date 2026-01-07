package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (r *Repository) Create(ctx context.Context, p *post.Post) error {
	query := `
		INSERT INTO posts (
			id,
			title,
			slug,
			content,
			author_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		p.ID,
		p.Title,
		p.Slug,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return err
}
