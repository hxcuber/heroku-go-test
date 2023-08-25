package user

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
)

func (i impl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := orm.Users(orm.UserWhere.UserEmail.EQ(email)).Bind(ctx, i.dbConn, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrEmailNotFound
		}
		return model.User{}, err
	}
	return user, nil
}
