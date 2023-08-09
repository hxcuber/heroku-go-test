package twoEmails

import (
	"errors"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	switch {
	case len(req.Friends) == 0:
		return errors.New("friends is a required field\n")
	case len(req.Friends) < 2:
		return errors.New("2 elements required, less than 2 given\n")
	case len(req.Friends) > 2:
		return errors.New("2 elements required, more than 2 given\n")
	case len(req.Friends) == 2:
		for _, e := range req.Friends {
			if err := request.IsEmailError(e, "elements of friends"); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("unknown error\n")
}
