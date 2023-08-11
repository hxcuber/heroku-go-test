package email

import (
	"errors"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	if req.Email == "" {
		return errors.New("email is a required field\n")
	}
	return request.IsEmailError(req.Email, "email")
}
