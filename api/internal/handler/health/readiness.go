package health

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/handler"
	"github.com/pkg/errors"
	"net/http"
)

// CheckReadiness checks for system readiness
func (h Handler) CheckReadiness() http.HandlerFunc {
	return handler.ErrorHandler(func(w http.ResponseWriter, r *http.Request) (error, int) {
		err := h.systemCtrl.CheckReadiness(r.Context())

		if errors.Is(err, context.Canceled) {
			return nil, http.StatusOK
		}

		return err, http.StatusInternalServerError
	})
}
