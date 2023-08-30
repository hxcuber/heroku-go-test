package relationship

import (
	"github.com/pkg/errors"
)

var (
	ErrAlreadyCreated = errors.New("relationship already created")
	ErrBlocked        = errors.New("one user has already blocked another")
	ErrFriends        = errors.New("users are already friends")
)
