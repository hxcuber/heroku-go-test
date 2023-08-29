package handler

import (
	"github.com/go-chi/render"
	"github.com/hxcuber/friends-management/api/internal/handler/response/errorWithString"
	"net/http"
)

func ErrorHandler(f func(w http.ResponseWriter, r *http.Request) (error, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err, status := f(w, r)

		if err != nil {
			render.Render(w, r, errorWithString.New(err.Error(), status))
			return
		}
	}
}
