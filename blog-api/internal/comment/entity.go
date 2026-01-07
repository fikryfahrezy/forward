package comment

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID    `json:"id"`
	Content   string       `json:"content"`
	PostID    uuid.UUID    `json:"post_id"`
	AuthorID  uuid.UUID    `json:"author_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

type CommentWithAuthor struct {
	Comment
	AuthorUsername string `json:"author_username"`
}

type CommentID struct {
	ID string `json:"id"`
}

type CreateCommentRequest struct {
	Content string `json:"content" example:"Great post! Thanks for sharing."`
}

func (r CreateCommentRequest) Validate() error {
	if r.Content == "" {
		return ErrInvalidInput
	}
	return nil
}

type UpdateCommentRequest struct {
	Content string `json:"content" example:"Updated comment content."`
}

func (r UpdateCommentRequest) Validate() error {
	if r.Content == "" {
		return ErrInvalidInput
	}
	return nil
}

type CommentItem struct {
	ID             uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content        string    `json:"content" example:"Great post! Thanks for sharing."`
	PostID         uuid.UUID `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	AuthorID       uuid.UUID `json:"author_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	AuthorUsername string    `json:"author_username" example:"johndoe"`
	CreatedAt      time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type CommentListResponse struct {
	Comments   []CommentItem `json:"comments"`
	TotalCount int           `json:"total_count" example:"50"`
	Page       int           `json:"page" example:"1"`
	PageSize   int           `json:"page_size" example:"10"`
}

func (c *CommentWithAuthor) ToCommentItem() CommentItem {
	return CommentItem{
		ID:             c.ID,
		Content:        c.Content,
		PostID:         c.PostID,
		AuthorID:       c.AuthorID,
		AuthorUsername: c.AuthorUsername,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
