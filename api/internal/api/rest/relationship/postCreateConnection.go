package relationship

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/twoEmails"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"net/http"
)

func (h Handler) PostCreateConnection() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request twoEmails.Request
		err := render.Bind(r, &request)
		if err != nil {
			return err, http.StatusBadRequest
		}

		err = h.ctrl.CreateConnection(r.Context(), request.Friends[0], request.Friends[1])
		if err != nil {
			if errors.Is(err, controller.ErrAlreadyCreated) {
				return err, http.StatusConflict
			}
			return err, http.StatusInternalServerError
		}

		err = render.Render(w, r, basicSuccess.New(http.StatusCreated))
		if err != nil {
			return err, http.StatusInternalServerError
		}

		return nil, http.StatusCreated
	})
}
