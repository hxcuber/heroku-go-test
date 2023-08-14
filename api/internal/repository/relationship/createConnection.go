package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) CreateConnection(ctx context.Context, user1 model.User, user2 model.User) (bool, error) {
	rela1to2 := &orm.Relationship{
		SenderID:   user1.UserID,
		ReceiverID: user2.UserID,
		Friends:    true,
	}

	rela2to1 := &orm.Relationship{
		SenderID:   user2.UserID,
		ReceiverID: user1.UserID,
		Friends:    true,
	}

	conflictColumns := []string{orm.RelationshipColumns.SenderID, orm.RelationshipColumns.ReceiverID}
	updateColumns := boil.Whitelist(orm.RelationshipColumns.Friends)
	insertColumns := boil.Whitelist(orm.RelationshipColumns.SenderID, orm.RelationshipColumns.ReceiverID, orm.RelationshipColumns.Friends)

	err := rela1to2.Upsert(ctx, i.dbConn, true, conflictColumns, updateColumns, insertColumns)
	if err != nil {
		return false, err
	}

	err = rela2to1.Upsert(ctx, i.dbConn, true, conflictColumns, updateColumns, insertColumns)
	if err != nil {
		return false, err
	}

	return rela1to2.Status != "" || rela2to1.Status != "", nil
}
