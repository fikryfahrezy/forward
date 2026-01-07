package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) Delete(ctx context.Context, postID, authorID uuid.UUID) (post.PostID, error) {
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

	if err := s.repo.Delete(ctx, postID); err != nil {
		return post.PostID{}, err
	}

	return post.PostID{ID: postID.String()}, nil
}
