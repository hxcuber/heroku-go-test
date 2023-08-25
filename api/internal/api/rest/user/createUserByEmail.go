package user

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/email"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"net/http"
)

func (h Handler) CreateUserByEmail() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request email.Request
		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		if err := h.ctrl.CreateUserByEmail(r.Context(), request.Email); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		if err := render.Render(w, r, basicSuccess.New(http.StatusCreated)); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		return nil, http.StatusCreated
	})
}
