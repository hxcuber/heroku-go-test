package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/senderText"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/recipients"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) Receivers() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request senderText.Request

		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		receivers, err := h.ctrl.Receivers(r.Context(), request.Sender, request.Text)
		if err != nil {
			if errors.Is(err, user.ErrEmailNotFound) {
				return err, http.StatusNotFound
			}
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		if err = render.Render(w, r, recipients.New(receivers, http.StatusOK)); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		return nil, http.StatusOK
	})
}
