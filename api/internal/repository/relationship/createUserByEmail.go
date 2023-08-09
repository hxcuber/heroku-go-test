package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) CreateUserByEmail(ctx context.Context, email string) error {
	newUser := &orm.User{
		UserID:    0,
		UserEmail: email,
	}
	insertColumns := boil.Blacklist(orm.UserColumns.UserID)

	if err := newUser.Insert(ctx, i.dbConn, insertColumns); err != nil {
		return err
	}

	return nil
}
