package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrInvalidUser       = errors.New("invalid user")
	ErrInvalidToken      = errors.New("invalid session token")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidCSRFToken  = errors.New("invalid csrf token")
	ErrNoRecord          = errors.New("record not found")
)
