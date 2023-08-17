package email

import (
	"github.com/hxcuber/friends-management/api/internal/api/rest/request"
	"github.com/pkg/errors"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	if req.Email == "" {
		return errors.New("email is a required field\n")
	}
	return request.IsEmailError(req.Email, "email")
}
