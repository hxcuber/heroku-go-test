package user

import (
	"github.com/pkg/errors"
)

var (
	ErrAlreadyCreated = errors.New("user already created")
)
