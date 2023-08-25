package user

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

type Controller interface {
	CreateUserByEmail(ctx context.Context, email string) error
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
