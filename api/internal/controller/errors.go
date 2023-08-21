package controller

import "github.com/pkg/errors"

var (
	ErrAlreadyCreated = errors.New("resource already created")
	ErrNotFound       = errors.New("resource not found")
	ErrBlocked        = errors.New("one user has already blocked another")
	ErrFriends        = errors.New("users are already friends")
)
