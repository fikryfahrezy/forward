package post

import (
	appError "github.com/fikryfahrezy/forward/blog-api/internal/error"
)

var (
	ErrPostNotFound   = appError.New("POST_NOT_FOUND", "Post not found")
	ErrInvalidInput   = appError.New("INVALID_INPUT", "Invalid input data")
	ErrUnauthorized   = appError.New("UNAUTHORIZED", "You are not authorized to perform this action")
	ErrSlugExists     = appError.New("SLUG_EXISTS", "A post with this slug already exists")
)
