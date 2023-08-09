package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
)

func (i impl) GetNotificationStatus(ctx context.Context, senderEmail string, receiverEmail string) (string, error) {
	sender, err := i.getUserByEmail(ctx, senderEmail)
	if err != nil {
		return "", err
	}

	receiver, err := i.getUserByEmail(ctx, receiverEmail)
	if err != nil {
		return "", err
	}

	relationship, err := orm.FindRelationship(ctx, i.dbConn, sender.UserID, receiver.UserID, orm.RelationshipColumns.Status)
	if err != nil {
		return "", err
	}

	return relationship.Status, nil
}
