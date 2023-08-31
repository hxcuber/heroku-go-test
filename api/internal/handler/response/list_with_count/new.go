package list_with_count

import (
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
)

type Response struct {
	basic_success.Response
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func New(list []string, statusCode int) *Response {
	return &Response{
		Response: basic_success.Response{
			Success:        true,
			HttpStatusCode: statusCode,
		},
		Friends: list,
		Count:   len(list),
	}
}
