package relationship

import (
	"github.com/hxcuber/friends-management/api/internal/repository"
)

type Controller interface {
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
