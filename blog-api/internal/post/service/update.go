package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) Update(ctx context.Context, postID, authorID uuid.UUID, req post.UpdatePostRequest) (post.PostID, error) {
	if err := req.Validate(); err != nil {
		return post.PostID{}, err
	}

	p, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return post.PostID{}, err
	}
	if p == (post.Post{}) {
		return post.PostID{}, post.ErrPostNotFound
	}

	// Check if the user is the author
	if p.AuthorID != authorID {
		return post.PostID{}, post.ErrUnauthorized
	}

	// Generate new slug if title changed
	newSlug := p.Slug
	if p.Title != req.Title {
		newSlug, err = generateSlug(req.Title)
		if err != nil {
			return post.PostID{}, err
		}
		// Ensure new slug is unique (excluding current post)
		for {
			exists, err := s.repo.ExistsBySlugExcludingID(ctx, newSlug, postID)
			if err != nil {
				return post.PostID{}, err
			}
			if !exists {
				break
			}
			newSlug, err = generateSlug(req.Title)
			if err != nil {
				return post.PostID{}, err
			}
		}
	}

	p.Title = req.Title
	p.Slug = newSlug
	p.Content = req.Content
	p.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, p); err != nil {
		return post.PostID{}, err
	}

	return post.PostID{ID: p.ID.String()}, nil
}
