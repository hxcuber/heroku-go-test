package senderText

import (
	"errors"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	if req.Sender == "" {
		return errors.New("sender is a required field\n")
	}
	return request.IsEmail(req.Sender, "sender")
}
