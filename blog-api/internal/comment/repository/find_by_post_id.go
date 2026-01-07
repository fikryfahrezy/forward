package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (r *Repository) FindByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]comment.CommentWithAuthor, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT
			c.id,
			c.content,
			c.post_id,
			c.author_id,
			c.created_at,
			c.updated_at,
			u.username,
			COUNT(*) OVER() AS total_count
		FROM comments c
			JOIN users u ON c.author_id = u.id
		WHERE
			c.post_id = $1
			AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, postID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var totalCount int
	var comments []comment.CommentWithAuthor
	for rows.Next() {
		var c comment.CommentWithAuthor
		if err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.PostID,
			&c.AuthorID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.AuthorUsername,
			&totalCount,
		); err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}

	return comments, totalCount, nil
}
