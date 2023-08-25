package user

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUserByEmail(ctx context.Context, email string) error
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
