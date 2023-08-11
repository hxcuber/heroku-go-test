package request

import (
	"errors"
	"fmt"
	"github.com/hxcuber/friends-management/api/pkg/util"
)

func IsEmailError(s string, fieldName string) error {
	if !util.IsEmail(s) {
		return errors.New(fmt.Sprintf("%s must be an email\n", fieldName))
	}
	return nil
}
