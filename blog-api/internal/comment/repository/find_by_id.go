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

func (r *Repository) FindByIDWithAuthor(ctx context.Context, id uuid.UUID) (comment.CommentWithAuthor, error) {
	query := `
		SELECT
			c.id,
			c.content,
			c.post_id,
			c.author_id,
			c.created_at,
			c.updated_at,
			u.username
		FROM comments c
			JOIN users u ON c.author_id = u.id
		WHERE
			c.id = $1
			AND c.deleted_at IS NULL
	`
	c := comment.CommentWithAuthor{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Content,
		&c.PostID,
		&c.AuthorID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.AuthorUsername,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return comment.CommentWithAuthor{}, nil
	}
	if err != nil {
		return comment.CommentWithAuthor{}, err
	}
	return c, nil
}
