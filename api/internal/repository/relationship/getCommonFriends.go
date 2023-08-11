package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i impl) GetCommonFriends(ctx context.Context, user1 model.User, user2 model.User) (model.UserSlice, error) {
	// See https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices
	var friendList model.UserSlice
	err := orm.Users(
		qm.InnerJoin(fmt.Sprintf("%s on %s=%s",
			orm.TableNames.Relationships,
			orm.UserTableColumns.UserID,
			orm.RelationshipTableColumns.ReceiverID)),
		qm.Expr(
			orm.RelationshipWhere.SenderID.EQ(user1.UserID),
			qm.Or2(orm.RelationshipWhere.SenderID.EQ(user2.UserID)),
		),
		orm.RelationshipWhere.Friends.EQ(true),
	).Bind(ctx, i.dbConn, &friendList)
	if err != nil {
		return nil, err
	}
	return friendList, nil
}
