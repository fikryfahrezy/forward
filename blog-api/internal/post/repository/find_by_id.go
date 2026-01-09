package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (post.Post, error) {
	query := `
		SELECT
			id,
			title,
			slug,
			content,
			author_id,
			created_at,
			updated_at
		FROM posts
		WHERE 
			id = $1 
			AND deleted_at IS NULL
	`
	p := post.Post{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.AuthorID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return post.Post{}, nil
	}
	if err != nil {
		return post.Post{}, err
	}
	return p, nil
}
