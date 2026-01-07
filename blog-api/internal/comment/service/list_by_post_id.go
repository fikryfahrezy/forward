package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
)

func (s *Service) ListByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) (comment.CommentListResponse, error) {
	// Check if post exists
	exists, err := s.repo.PostExists(ctx, postID)
	if err != nil {
		return comment.CommentListResponse{}, err
	}
	if !exists {
		return comment.CommentListResponse{}, comment.ErrPostNotFound
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	comments, totalCount, err := s.repo.FindByPostID(ctx, postID, page, pageSize)
	if err != nil {
		return comment.CommentListResponse{}, err
	}

	commentItems := make([]comment.CommentItem, len(comments))
	for i, c := range comments {
		commentItems[i] = c.ToCommentItem()
	}

	return comment.CommentListResponse{
		Comments:   commentItems,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
