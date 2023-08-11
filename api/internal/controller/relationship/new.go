package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

type Controller interface {
	GetFriends(ctx context.Context, email string) ([]string, error)
	GetCommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error)
	GetReceivers(ctx context.Context, email string, text string) ([]string, error)
	CreateUserByEmail(ctx context.Context, email string) error
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
