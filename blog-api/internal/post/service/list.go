package service

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) List(ctx context.Context, page, pageSize int) (post.PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, totalCount, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return post.PostListResponse{}, err
	}

	postItems := make([]post.PostItem, len(posts))
	for i, p := range posts {
		postItems[i] = p.ToPostItem()
	}

	return post.PostListResponse{
		Posts:      postItems,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
