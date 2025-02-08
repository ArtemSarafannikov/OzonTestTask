package cstErrors

import "errors"

var (
	UserAlreadyExistsError = GenerateError("User already exists")
	UnauthorizedError      = GenerateError("Authorize for this action")
	InternalError          = GenerateError("Internal server error")
	TooLongContentError    = GenerateError("Too long content")
	NotFoundError          = GenerateError("Not found")
	InvalidCredentials     = GenerateError("Invalid credentials")
	PermissionDeniedError  = GenerateError("Permission denied")
	CommentNotInPostError  = GenerateError("Comment not in post")
)

func GenerateError(err string) error {
	return errors.New(err)
}
