package cstErrors

import "errors"

type KnownError interface {
	IsKnown() bool
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func (e *customError) IsKnown() bool {
	return true
}

var (
	UserAlreadyExistsError = GenerateError("User already exists")
	UnauthorizedError      = GenerateError("Authorize for this action")
	InternalError          = GenerateError("Internal server error")
	TooLongContentError    = GenerateError("Too long content")
	NotFoundError          = GenerateError("Not found")
	InvalidCredentials     = GenerateError("Invalid credentials")
	PermissionDeniedError  = GenerateError("Permission denied")
	CommentNotInPostError  = GenerateError("Comment not in post")
	NoJWTSecretError       = GenerateError("No jwt secret in environment")
)

func GenerateError(err string) error {
	return &customError{msg: err}
}

func IsCustomError(err error) bool {
	var knownError KnownError
	return errors.As(err, &knownError)
}

func GetCustomError(err error) error {
	if err == nil {
		return nil
	}
	if !IsCustomError(err) {
		return InternalError
	}
	return err
}
