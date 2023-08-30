package user

import (
	"github.com/hxcuber/friends-management/api/internal/controller/user"
)

type Handler struct {
	ctrl user.Controller
}

func New(ctrl user.Controller) Handler {
	return Handler{ctrl: ctrl}
}
