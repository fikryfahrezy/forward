package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) Create(ctx context.Context, authorID uuid.UUID, req post.CreatePostRequest) (post.PostID, error) {
	if err := req.Validate(); err != nil {
		return post.PostID{}, err
	}

	slug, err := generateSlug(req.Title)
	if err != nil {
		return post.PostID{}, err
	}

	// Ensure slug is unique
	for {
		exists, err := s.repo.ExistsBySlug(ctx, slug)
		if err != nil {
			return post.PostID{}, err
		}
		if !exists {
			break
		}
		slug, err = generateSlug(req.Title) // Regenerate with new random suffix
		if err != nil {
			return post.PostID{}, err
		}
	}

	now := time.Now()
	p := &post.Post{
		ID:        uuid.Must(uuid.NewV7()),
		Title:     req.Title,
		Slug:      slug,
		Content:   req.Content,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return post.PostID{}, err
	}

	return post.PostID{ID: p.ID.String()}, nil
}
