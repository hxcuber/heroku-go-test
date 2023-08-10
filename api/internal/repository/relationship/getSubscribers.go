package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i impl) GetSubscribers(ctx context.Context, senderEmail string) (model.UserSlice, error) {
	sender, err := i.getUserByEmail(ctx, senderEmail)
	if err != nil {
		return nil, err
	}

	var subscriberList model.UserSlice
	err = orm.Users(
		qm.InnerJoin(fmt.Sprintf("%s on %s=%s",
			orm.TableNames.Relationships,
			orm.UserTableColumns.UserID,
			orm.RelationshipTableColumns.ReceiverID)),
		orm.RelationshipWhere.SenderID.EQ(sender.UserID),
		qm.Expr(
			orm.RelationshipWhere.Status.EQ(orm.SubscriptionStatusRSubscribedS),
			qm.Or2(
				qm.Expr(
					orm.RelationshipWhere.Friends.EQ(true),
					orm.RelationshipWhere.Status.EQ(orm.SubscriptionStatusNone)))),
	).Bind(ctx, i.dbConn, &subscriberList)
	if err != nil {
		return nil, err
	}
	return subscriberList, nil
}
