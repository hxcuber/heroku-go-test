package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/requestorTarget"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) Subscribe() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request requestorTarget.Request
		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		err := h.ctrl.Subscribe(r.Context(), request.Requestor, request.Target)
		if err != nil {
			if errors.Is(err, user.ErrEmailNotFound) {
				return err, http.StatusNotFound
			}
			if errors.Is(err, relationship.ErrAlreadyCreated) {
				return errors.Wrap(err, "subscription"), http.StatusConflict
			}
			return err, http.StatusInternalServerError
		}

		if err = render.Render(w, r, basicSuccess.New(http.StatusOK)); err != nil {
			return err, http.StatusInternalServerError
		}

		return nil, http.StatusOK
	})
}
