package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/twoEmails"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/listWithCount"
	"net/http"
)

func (h Handler) GetCommonFriends() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request twoEmails.Request

		if err := render.Bind(r, &request); err != nil {
			return err, http.StatusBadRequest
		}

		list, err := h.ctrl.GetCommonFriends(r.Context(), request.Friends[0], request.Friends[1])
		if err != nil {
			return err, http.StatusInternalServerError
		}

		// purely for JSON conversion purposes
		if list == nil {
			list = []string{}
		}

		if err = render.Render(w, r, listWithCount.New(list, http.StatusOK)); err != nil {
			return err, http.StatusInternalServerError
		}

		return nil, http.StatusOK
	})
}
