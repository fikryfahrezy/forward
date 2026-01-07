package user

import (
	appError "github.com/fikryfahrezy/forward/blog-api/internal/error"
)

var (
	ErrUserNotFound       = appError.New("USER_NOT_FOUND", "User not found")
	ErrInvalidCredentials = appError.New("INVALID_CREDENTIALS", "Invalid email or password")
	ErrEmailExists        = appError.New("EMAIL_EXISTS", "Email already exists")
	ErrUsernameExists     = appError.New("USERNAME_EXISTS", "Username already exists")
	ErrInvalidInput       = appError.New("INVALID_INPUT", "Invalid input data")
)
