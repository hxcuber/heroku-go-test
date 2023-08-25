package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
)

type Repository interface {
	GetFriends(ctx context.Context, user model.User) (model.Users, error)
	GetSubscribers(ctx context.Context, sender model.User) (model.Users, error)
	GetReceiversFromEmails(ctx context.Context, sender model.User, emails []string) (model.Users, error)
	CreateRelationship(ctx context.Context, rela model.Relationship) error
	FindRelationship(ctx context.Context, receiver model.User, sender model.User) (*model.Relationship, error)
	UpdateRelationship(ctx context.Context, rela model.Relationship) error
	DeleteRelationship(ctx context.Context, rela model.Relationship) error
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
