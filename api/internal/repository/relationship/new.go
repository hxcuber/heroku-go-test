package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
)

type Repository interface {
	GetFriends(ctx context.Context, email string) (model.UserSlice, error)
	GetSubscribers(ctx context.Context, email string) (model.UserSlice, error)
	GetNotificationStatus(ctx context.Context, senderEmail string, receiverEmail string) (string, error)
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
