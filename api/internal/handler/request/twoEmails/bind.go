package twoEmails

import (
	"github.com/hxcuber/friends-management/api/internal/handler/request"
	"github.com/pkg/errors"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	switch {
	case len(req.Friends) == 0:
		return errors.New("friends is a required field")
	case len(req.Friends) < 2:
		return errors.New("2 elements required, less than 2 given")
	case len(req.Friends) > 2:
		return errors.New("2 elements required, more than 2 given")
	case len(req.Friends) == 2:
		for _, e := range req.Friends {
			if err := request.IsEmailError(e, "elements of friends"); err != nil {
				return err
			}
		}
		if req.Friends[0] == req.Friends[1] {
			return errors.New("emails cannot be the same")
		}
		return nil
	}
	return errors.New("unknown error")
}
