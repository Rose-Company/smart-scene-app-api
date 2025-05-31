package common

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrVideoNotFound    = errors.New("video not found")
	ErrInvalidUUID      = errors.New("invalid UUID format")
) 