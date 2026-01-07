package repository

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]post.PostWithAuthor, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT
			p.id,
			p.title,
			p.slug,
			p.content,
			p.author_id,
			p.created_at,
			p.updated_at,
			u.username,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE
			p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var totalCount int
	var posts []post.PostWithAuthor
	for rows.Next() {
		var p post.PostWithAuthor
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.AuthorID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.AuthorUsername,
			&totalCount,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	return posts, totalCount, nil
}
