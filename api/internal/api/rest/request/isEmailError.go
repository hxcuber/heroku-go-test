package request

import (
	"fmt"
	"github.com/hxcuber/friends-management/api/pkg/util"
	"github.com/pkg/errors"
)

func IsEmailError(s string, fieldName string) error {
	if !util.IsEmail(s) {
		return errors.New(fmt.Sprintf("%s must be an email", fieldName))
	}
	return nil
}
