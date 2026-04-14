package db_errors

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrDuplicateUser = errors.New("user already exists")
)

var (
	ErrFlightNotFound = errors.New("flight not found")
)
