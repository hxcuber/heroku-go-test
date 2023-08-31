package user

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/controller/user"
	"github.com/hxcuber/friends-management/api/internal/handler"
	"github.com/hxcuber/friends-management/api/internal/handler/request/email"
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
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
			if errors.Is(err, user.ErrAlreadyCreated) {
				return err, http.StatusConflict
			}
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		if err := render.Render(w, r, basic_success.New(http.StatusCreated)); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		return nil, http.StatusCreated
	})
}
