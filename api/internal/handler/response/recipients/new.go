package recipients

import (
	"github.com/hxcuber/friends-management/api/internal/handler/response/basicSuccess"
)

type Response struct {
	basicSuccess.Response
	Recipients []string `json:"recipients"`
}

func New(list []string, statusCode int) *Response {
	return &Response{
		Response: basicSuccess.Response{
			Success:        true,
			HttpStatusCode: statusCode,
		},
		Recipients: list,
	}
}
