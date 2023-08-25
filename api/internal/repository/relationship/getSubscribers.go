package relationship

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i impl) GetSubscribers(ctx context.Context, sender model.User) (model.Users, error) {
	var subscriberList model.Users
	err := orm.Users(
		qm.InnerJoin(fmt.Sprintf("%s on %s=%s",
			orm.TableNames.Relationships,
			orm.UserTableColumns.UserID,
			orm.RelationshipTableColumns.ReceiverID)),
		orm.RelationshipWhere.SenderID.EQ(sender.UserID),
		orm.RelationshipWhere.Status.NEQ(orm.StatusRBlockedS),
	).Bind(ctx, i.dbConn, &subscriberList)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return subscriberList, nil
}
