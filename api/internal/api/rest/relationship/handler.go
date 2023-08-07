package relationship

import "github.com/hxcuber/friends-management/api/internal/controller/relationship"

type Handler struct {
	ctrl relationship.Controller
}

func New(ctrl relationship.Controller) Handler {
	return Handler{ctrl: ctrl}
}
