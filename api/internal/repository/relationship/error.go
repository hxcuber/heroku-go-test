package relationship

import "github.com/pkg/errors"

var (
	ErrEmailNotFound = errors.New("email not found in database")
)
