package common

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
) 