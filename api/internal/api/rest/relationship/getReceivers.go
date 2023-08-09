package relationship

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/api/rest"
	"github.com/hxcuber/friends-management/api/internal/api/rest/request/senderText"
	"net/http"
)

func (h Handler) GetReceivers() http.HandlerFunc {
	return rest.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		var request senderText.Request
		err := render.Bind(r, &request)
		if err != nil {
			return err, http.StatusBadRequest
		}
		return nil, 0
	})
}
