package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/twoEmails"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/basicSuccess"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) Befriend() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request twoEmails.Request

		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		err := h.ctrl.Befriend(r.Context(), request.Friends[0], request.Friends[1])
		if err != nil {
			if errors.Is(err, user.ErrEmailNotFound) {
				return err, http.StatusNotFound
			}
			if errors.Is(err, relationship.ErrAlreadyCreated) {
				return errors.Wrap(err, "friendship"), http.StatusConflict
			}
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		if err = render.Render(w, r, basicSuccess.New(http.StatusCreated)); err != nil {
			return errors.New("Something went wrong"), http.StatusInternalServerError
		}

		return nil, http.StatusCreated
	})
}
