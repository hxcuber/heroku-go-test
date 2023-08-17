package health

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/hxcuber/friends-management/api/pkg/httpserv"
)

// CheckReadiness checks for system readiness
func (h Handler) CheckReadiness() http.HandlerFunc {
	return httpserv.ErrHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		err := h.systemCtrl.CheckReadiness(r.Context())

		if errors.Is(err, context.Canceled) {
			return nil
		}

		return err
	})
}
