package error

import (
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
)

type Response struct {
	basicSuccess.Response
	ErrMessage string `json:"error_message"`
}

func New(message string, statusCode int) Response {
	return Response{
		Response: basicSuccess.Response{
			Success:        false,
			HttpStatusCode: statusCode,
		},
		ErrMessage: message,
	}
}
