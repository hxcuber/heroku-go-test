package relationship

import (
	"github.com/pkg/errors"
)

var (
	ErrAlreadyCreated = errors.New("relationship already created")
	ErrNotFound       = errors.New("relationship not found")
	ErrBlocked        = errors.New("one user has already blocked another")
	ErrFriends        = errors.New("users are already friends")
)
