package relationship

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrAlreadyCreated = errors.New("resource already created")
	ErrNotFound       = errors.New("resource not found")
	ErrBlocked        = errors.New("one user has already blocked another")
	ErrFriends        = errors.New("users are already friends")
)

func LogErrMessage(controllerName string, message string, err error, values ...any) string {
	return fmt.Sprintf("[%s] %s encountered an error: %#v\n", controllerName, fmt.Sprintf(message, values...), err)
}
