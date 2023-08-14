package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
)

type Repository interface {
	GetFriends(ctx context.Context, user model.User) (model.UserSlice, error)
	GetSubscribers(ctx context.Context, sender model.User) (model.UserSlice, error)
	GetReceiversFromEmails(ctx context.Context, sender model.User, emails []string) (model.UserSlice, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateConnection(ctx context.Context, user1 model.User, user2 model.User) (bool, error)
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
