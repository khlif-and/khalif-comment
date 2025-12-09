package domain

import "errors"

var (
	// Generic Errors
	ErrInternalServerError = errors.New("internal server error")
	ErrBadParamInput       = errors.New("given param is not valid")

	// Resource Errors
	ErrCommentNotFound = errors.New("comment not found")
	ErrStoryNotFound   = errors.New("story not found") // Opsional, jika validasi story dilakukan di sini

	// Permission & Validation Errors
	ErrUnauthorizedAction = errors.New("you are not authorized to modify this comment") // User hanya boleh edit/delete punya sendiri
	ErrEmptyContent       = errors.New("comment content cannot be empty")
)