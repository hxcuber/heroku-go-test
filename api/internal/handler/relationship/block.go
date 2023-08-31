package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/handler"
	"github.com/hxcuber/friends-management/api/internal/handler/request/requestor_target"
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) Block() http.HandlerFunc {
	return handler.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request requestor_target.Request
		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		err := h.ctrl.Block(r.Context(), request.Requestor, request.Target)
		if err != nil {
			if errors.Is(err, user.ErrEmailNotFound) {
				return err, http.StatusNotFound
			}
			if errors.Is(err, relationship.ErrAlreadyCreated) || errors.Is(err, relationship.ErrBlocked) {
				return errors.Wrap(err, "block"), http.StatusConflict
			}
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		if err = render.Render(w, r, basic_success.New(http.StatusOK)); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		return nil, http.StatusOK
	})
}
