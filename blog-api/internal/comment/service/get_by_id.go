package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (comment.CommentWithAuthor, error) {
	c, err := s.repo.FindByIDWithAuthor(ctx, id)
	if err != nil {
		return comment.CommentWithAuthor{}, err
	}
	if c == (comment.CommentWithAuthor{}) {
		return comment.CommentWithAuthor{}, comment.ErrCommentNotFound
	}
	return c, nil
}
