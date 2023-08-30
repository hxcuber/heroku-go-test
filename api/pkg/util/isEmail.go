package util

import "strings"

func IsEmail(s string) bool {
	partsAt := strings.Split(s, "@")
	if len(partsAt) != 2 {
		return false
	}

	partsDot := strings.Split(partsAt[1], ".")
	if len(partsDot) != 2 {
		return false
	}
	if partsDot[1] == "" {
		return false
	}

	return true
}
