package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

type Controller interface {
	Friends(ctx context.Context, email string) ([]string, error)
	CommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error)
	Receivers(ctx context.Context, email string, text string) ([]string, error)
	Block(ctx context.Context, requestorEmail string, targetEmail string) error
	Befriend(ctx context.Context, email1 string, email2 string) error
	Subscribe(ctx context.Context, requestorEmail string, targetEmail string) error
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
