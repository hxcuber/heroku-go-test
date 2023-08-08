package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
)

func (i impl) getUserByEmail(ctx context.Context, email string) (*orm.User, error) {
	return orm.Users(orm.UserWhere.UserEmail.EQ(email)).One(ctx, i.dbConn)
}
