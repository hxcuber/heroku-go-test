package recipients

import (
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
)

type Response struct {
	basic_success.Response
	Recipients []string `json:"recipients"`
}

func New(list []string, statusCode int) *Response {
	return &Response{
		Response: basic_success.Response{
			Success:        true,
			HttpStatusCode: statusCode,
		},
		Recipients: list,
	}
}
