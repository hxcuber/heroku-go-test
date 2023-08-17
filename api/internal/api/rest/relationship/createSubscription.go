package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/requestorTarget"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) CreateSubscription() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request requestorTarget.Request
		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		err := h.ctrl.CreateSubscription(r.Context(), request.Requestor, request.Target)
		if err != nil {
			if errors.Is(err, controller.ErrAlreadyCreated) {
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
