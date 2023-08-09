package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/senderText"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/recipients"
	"net/http"
)

func (h Handler) GetReceivers() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request senderText.Request
		err := render.Bind(r, &request)
		if err != nil {
			return err, http.StatusBadRequest
		}

		receivers, err := h.ctrl.GetReceivers(r.Context(), request.Sender, request.Text)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		err = render.Render(w, r, recipients.New(receivers, http.StatusOK))
		if err != nil {
			return err, http.StatusInternalServerError
		}
		return nil, http.StatusOK
	})
}
