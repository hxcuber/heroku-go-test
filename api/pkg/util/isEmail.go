package util

import "strings"

func IsEmail(s string) bool {
	return strings.Index(s, "@") != -1
}
