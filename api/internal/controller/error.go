package controller

import "fmt"

func LogErrMessage(controllerName string, message string, err error, values ...any) string {
	return fmt.Sprintf("[%s] %s encountered an error: %#v\n", controllerName, fmt.Sprintf(message, values...), err)
}
