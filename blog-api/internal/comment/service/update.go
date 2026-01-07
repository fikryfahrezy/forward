package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (s *Service) Update(ctx context.Context, commentID, authorID uuid.UUID, req comment.UpdateCommentRequest) (comment.CommentID, error) {
	if err := req.Validate(); err != nil {
		return comment.CommentID{}, err
	}

	c, err := s.repo.FindByID(ctx, commentID)
	if err != nil {
		return comment.CommentID{}, err
	}
	if c == (comment.Comment{}) {
		return comment.CommentID{}, comment.ErrCommentNotFound
	}

	// Check if the user is the author
	if c.AuthorID != authorID {
		return comment.CommentID{}, comment.ErrUnauthorized
	}

	c.Content = req.Content
	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return comment.CommentID{}, err
	}

	return comment.CommentID{ID: c.ID.String()}, nil
}
