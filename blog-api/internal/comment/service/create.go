package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (s *Service) Create(ctx context.Context, postID, authorID uuid.UUID, req comment.CreateCommentRequest) (comment.CommentID, error) {
	if err := req.Validate(); err != nil {
		return comment.CommentID{}, err
	}

	// Check if post exists
	exists, err := s.repo.PostExists(ctx, postID)
	if err != nil {
		return comment.CommentID{}, err
	}
	if !exists {
		return comment.CommentID{}, comment.ErrPostNotFound
	}

	now := time.Now()
	c := &comment.Comment{
		ID:        uuid.Must(uuid.NewV7()),
		Content:   req.Content,
		PostID:    postID,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return comment.CommentID{}, err
	}

	return comment.CommentID{ID: c.ID.String()}, nil
}
