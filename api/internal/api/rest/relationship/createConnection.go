package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/twoEmails"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"net/http"
)

func (h Handler) CreateConnection() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request twoEmails.Request
		err := render.Bind(r, &request)
		if err != nil {
			return err, http.StatusBadRequest
		}

		created, err := h.ctrl.CreateConnection(r.Context(), request.Friends[0], request.Friends[1])
		if err != nil {
			return err, http.StatusInternalServerError
		}

		statusCode := http.StatusOK
		if created {
			statusCode = http.StatusNoContent
		}

		err = render.Render(w, r, basicSuccess.New(statusCode))
		if err != nil {
			return err, http.StatusInternalServerError
		}

		return nil, statusCode
	})
}
