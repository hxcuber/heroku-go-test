package senderText

import (
	"github.com/hxcuber/friends-management/api/internal/handler/request"
	"github.com/pkg/errors"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	if req.Sender == "" {
		return errors.New("sender is a required field")
	}
	return request.IsEmailError(req.Sender, "sender")
}
