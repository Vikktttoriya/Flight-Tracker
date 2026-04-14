package service_errors

import (
	"fmt"
)

type ErrorCode string

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

const (
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
	CodeAlreadyExists      ErrorCode = "ALREADY_EXISTS"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeInvalidArgument    ErrorCode = "INVALID_ARGUMENT"
	CodeForbidden          ErrorCode = "FORBIDDEN"
	CodeSelfModification   ErrorCode = "SELF_MODIFICATION"
	CodeInvalidTransition  ErrorCode = "INVALID_STATUS_TRANSITION"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
)
