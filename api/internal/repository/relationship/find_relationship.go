package relationship

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
)

func (i impl) FindRelationship(ctx context.Context, receiver model.User, sender model.User) (*model.Relationship, error) {
	relationshipOrm, err := orm.FindRelationship(ctx, i.dbConn, receiver.UserID, sender.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Relationship{}, ErrRelationshipNotFound
		}
		return &model.Relationship{}, err
	}
	return &model.Relationship{
		ReceiverID: relationshipOrm.ReceiverID,
		SenderID:   relationshipOrm.SenderID,
		Status:     relationshipOrm.Status,
	}, nil
}
