package requestorTarget

import (
	"errors"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request"
	"net/http"
)

func (req *Request) Bind(r *http.Request) error {
	if req.Requestor == "" {
		return errors.New("requestor is a required field")
	}

	if req.Target == "" {
		return errors.New("target is a required field")
	}

	if err := request.IsEmailError(req.Requestor, "requestor"); err != nil {
		return err
	}

	if err := request.IsEmailError(req.Target, "target"); err != nil {
		return err
	}

	return nil
}
