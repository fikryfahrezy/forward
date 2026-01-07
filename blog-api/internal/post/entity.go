package post

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Slug      string       `json:"slug"`
	Content   string       `json:"content"`
	AuthorID  uuid.UUID    `json:"author_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

type PostWithAuthor struct {
	Post
	AuthorUsername string `json:"author_username"`
}

type PostID struct {
	ID string `json:"id"`
}

type CreatePostRequest struct {
	Title   string `json:"title" example:"My First Blog Post"`
	Content string `json:"content" example:"This is the content of my first blog post..."`
}

func (r CreatePostRequest) Validate() error {
	if r.Title == "" || r.Content == "" {
		return ErrInvalidInput
	}
	return nil
}

type UpdatePostRequest struct {
	Title   string `json:"title" example:"Updated Blog Post Title"`
	Content string `json:"content" example:"This is the updated content..."`
}

func (r UpdatePostRequest) Validate() error {
	if r.Title == "" || r.Content == "" {
		return ErrInvalidInput
	}
	return nil
}

type PostItem struct {
	ID             uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string    `json:"title" example:"My First Blog Post"`
	Slug           string    `json:"slug" example:"my-first-blog-post-a1b2c3d4"`
	Content        string    `json:"content" example:"This is the content of my first blog post..."`
	AuthorID       uuid.UUID `json:"author_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	AuthorUsername string    `json:"author_username" example:"johndoe"`
	CreatedAt      time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type PostListResponse struct {
	Posts      []PostItem `json:"posts"`
	TotalCount int        `json:"total_count" example:"100"`
	Page       int        `json:"page" example:"1"`
	PageSize   int        `json:"page_size" example:"10"`
}

func (p *PostWithAuthor) ToPostItem() PostItem {
	return PostItem{
		ID:             p.ID,
		Title:          p.Title,
		Slug:           p.Slug,
		Content:        p.Content,
		AuthorID:       p.AuthorID,
		AuthorUsername: p.AuthorUsername,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
