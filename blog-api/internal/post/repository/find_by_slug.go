package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (r *Repository) FindBySlugWithAuthor(ctx context.Context, slug string) (post.PostWithAuthor, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.slug,
			p.content,
			p.author_id,
			p.created_at,
			p.updated_at,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE
			p.slug = $1
			AND p.deleted_at IS NULL
	`
	p := post.PostWithAuthor{}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.AuthorID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.AuthorUsername,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return post.PostWithAuthor{}, nil
	}
	if err != nil {
		return post.PostWithAuthor{}, err
	}
	return p, nil
}
