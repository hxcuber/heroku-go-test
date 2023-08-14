package relationship

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
)

func (i impl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := orm.Users(orm.UserWhere.UserEmail.EQ(email)).Bind(ctx, i.dbConn, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New(fmt.Sprintf("email %s not found in databse", email))
		}
		return model.User{}, err
	}
	return user, nil
}
