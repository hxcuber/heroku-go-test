package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
)

func (i impl) FindRelationship(ctx context.Context, user1 model.User, user2 model.User) (model.Relationship, error) {
	relationshipOrm, err := orm.FindRelationship(ctx, i.dbConn, user1.UserID, user2.UserID)
	if err != nil {
		return model.Relationship{}, err
	}
	return model.Relationship{
		SenderID:   relationshipOrm.SenderID,
		ReceiverID: relationshipOrm.ReceiverID,
		Friends:    relationshipOrm.Friends,
		Status:     relationshipOrm.Status,
	}, err
}
