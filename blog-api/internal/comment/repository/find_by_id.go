package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (comment.Comment, error) {
	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`
	c := comment.Comment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Content,
		&c.PostID,
		&c.AuthorID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return comment.Comment{}, nil
	}
	if err != nil {
		return comment.Comment{}, err
	}
	return c, nil
}
