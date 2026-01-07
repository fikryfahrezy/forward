package comment

import (
	appError "github.com/fikryfahrezy/forward/blog-api/internal/error"
)

var (
	ErrCommentNotFound = appError.New("COMMENT_NOT_FOUND", "Comment not found")
	ErrPostNotFound    = appError.New("POST_NOT_FOUND", "Post not found")
	ErrInvalidInput    = appError.New("INVALID_INPUT", "Invalid input data")
	ErrUnauthorized    = appError.New("UNAUTHORIZED", "You are not authorized to perform this action")
)
