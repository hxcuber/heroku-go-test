package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/twoEmails"
	"github.com/hxcuber/friends-management/api/internal/api/rest/response/listWithCount"
	"net/http"
)

func (h Handler) GetCommonFriendList() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request twoEmails.Request
		err := render.Bind(r, &request)
		if err != nil {
			return err, http.StatusBadRequest
		}

		list, err := h.ctrl.GetCommonFriendList(r.Context(), request.Friends[0], request.Friends[1])
		if err != nil {
			return err, http.StatusInternalServerError
		}

		err = render.Render(w, r, listWithCount.New(list, http.StatusOK))
		if err != nil {
			return err, http.StatusInternalServerError
		}
		return nil, http.StatusOK
	})
}
