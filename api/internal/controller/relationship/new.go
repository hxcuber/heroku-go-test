package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

type Controller interface {
	GetFriendList(ctx context.Context, email string) ([]string, error)
	GetCommonFriendList(ctx context.Context, email1 string, email2 string) ([]string, error)
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
