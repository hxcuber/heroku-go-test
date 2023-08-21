package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
)

func (i impl) FindRelationship(ctx context.Context, receiver model.User, sender model.User) (*model.Relationship, error) {
	relationshipOrm, err := orm.FindRelationship(ctx, i.dbConn, receiver.UserID, sender.UserID)
	if err != nil {
		return nil, err
	}
	return &model.Relationship{
		ReceiverID: relationshipOrm.ReceiverID,
		SenderID:   relationshipOrm.SenderID,
		Status:     relationshipOrm.Status,
	}, err
}
