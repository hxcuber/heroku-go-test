package user

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/handler"
	"github.com/hxcuber/friends-management/api/internal/handler/request/email"
	"github.com/hxcuber/friends-management/api/internal/handler/response/basicSuccess"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) CreateUserByEmail() http.HandlerFunc {
	return handler.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
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
