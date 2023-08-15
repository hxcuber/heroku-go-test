package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) CreateBlock(ctx context.Context, sender model.User, receiver model.User) error {
	rela := &orm.Relationship{
		SenderID:   sender.UserID,
		ReceiverID: receiver.UserID,
		Status:     orm.SubscriptionStatusRBlockedS,
	}

	conflictColumns := []string{orm.RelationshipColumns.SenderID, orm.RelationshipColumns.ReceiverID}
	updateColumns := boil.Whitelist(orm.RelationshipColumns.Status)
	insertColumns := boil.Whitelist(orm.RelationshipColumns.SenderID, orm.RelationshipColumns.ReceiverID, orm.RelationshipColumns.Status)

	return rela.Upsert(ctx, i.dbConn, true, conflictColumns, updateColumns, insertColumns)
}
