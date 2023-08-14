package relationship

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i impl) GetFriends(ctx context.Context, user model.User) (model.UserSlice, error) {
	// See https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices
	var friendList model.UserSlice
	err := orm.Users(
		qm.InnerJoin(fmt.Sprintf("%s on %s=%s",
			orm.TableNames.Relationships,
			orm.UserTableColumns.UserID,
			orm.RelationshipTableColumns.ReceiverID)),
		orm.RelationshipWhere.SenderID.EQ(user.UserID),
		orm.RelationshipWhere.Friends.EQ(true),
	).Bind(ctx, i.dbConn, &friendList)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return friendList, nil
}
