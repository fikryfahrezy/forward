package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (s *Service) Delete(ctx context.Context, commentID, authorID uuid.UUID) (comment.CommentID, error) {
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

	if err := s.repo.Delete(ctx, commentID); err != nil {
		return comment.CommentID{}, err
	}

	return comment.CommentID{ID: commentID.String()}, nil
}
