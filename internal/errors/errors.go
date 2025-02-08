package cstErrors

import "errors"

var (
	UserAlreadyExistsError = GenerateError("User already exists")
	UnauthorizedError      = GenerateError("Authorize for this action")
	InternalError          = GenerateError("Internal server error")
	TooLongContentError    = GenerateError("Too long content")
	UserNotFoundError      = GenerateError("User not found")
	InvalidCredentials     = GenerateError("Invalid credentials")
)

func GenerateError(err string) error {
	return errors.New(err)
}
