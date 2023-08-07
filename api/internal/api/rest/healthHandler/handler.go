package healthHandler

import "github.com/hxcuber/friends-management/api/internal/controller/systemController"

// Handler is the web handler for this pkg
type Handler struct {
	systemCtrl systemController.Controller
}

// New instantiates a new Handler and returns it
func New(systemCtrl systemController.Controller) Handler {
	return Handler{systemCtrl: systemCtrl}
}
