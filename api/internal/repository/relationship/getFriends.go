package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i impl) GetFriends(ctx context.Context, email string) (model.UserSlice, error) {
	user, err := i.getUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var friendList model.UserSlice
	err = orm.Users(qm.InnerJoin(fmt.Sprintf("%s on %s=%s",
		orm.TableNames.Relationships,
		orm.UserTableColumns.UserID,
		orm.RelationshipTableColumns.ReceiverID)),
		orm.RelationshipWhere.SenderID.EQ(user.UserID),
		orm.RelationshipWhere.Friends.EQ(true),
	).Bind(ctx, i.dbConn, &friendList)
	if err != nil {
		return nil, err
	}
	return friendList, nil
}