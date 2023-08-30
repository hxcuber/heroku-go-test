package listWithCount

import (
	"github.com/hxcuber/friends-management/api/internal/handler/response/basicSuccess"
)

type Response struct {
	basicSuccess.Response
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func New(list []string, statusCode int) *Response {
	return &Response{
		Response: basicSuccess.Response{
			Success:        true,
			HttpStatusCode: statusCode,
		},
		Friends: list,
		Count:   len(list),
	}
}
