package request

import (
	"errors"
	"fmt"
	"strings"
)

func IsEmail(s string, fieldName string) error {
	if strings.Index(s, "@") == -1 {
		return errors.New(fmt.Sprintf("%s must be an email\n", fieldName))
	}
	return nil
}
