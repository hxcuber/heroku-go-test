package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
)

type Repository interface {
	GetFriendList(ctx context.Context, email string) (model.UserSlice, error)
	GetSubscribers(ctx context.Context, email string) (model.UserSlice, error)
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
