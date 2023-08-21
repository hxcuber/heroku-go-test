package relationship

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
)

func (i impl) Subscribe(ctx context.Context, requestorEmail string, targetEmail string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		receiver, err := txRepo.Relationship().GetUserByEmail(ctx, requestorEmail)
		if err != nil {
			return err
		}

		sender, err := txRepo.Relationship().GetUserByEmail(ctx, targetEmail)
		if err != nil {
			return err
		}

		relaStoR, err := txRepo.Relationship().FindRelationship(ctx, sender, receiver)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			if relaStoR.Status == orm.StatusRBlockedS {
				return controller.ErrBlocked
			}
		}

		relaRtoS, err := txRepo.Relationship().FindRelationship(ctx, receiver, sender)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			relaRtoS = &model.Relationship{
				ReceiverID: receiver.UserID,
				SenderID:   sender.UserID,
				Status:     orm.StatusRSubscribedS,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *relaRtoS)
		} else {
			switch relaRtoS.Status {
			case orm.StatusRSubscribedS:
				{
					return controller.ErrAlreadyCreated
				}
			case orm.StatusFriends:
				{
					return controller.ErrFriends
				}
			default:
				{
					relaRtoS.Status = orm.StatusRSubscribedS
				}
			}
			_, err = txRepo.Relationship().UpdateRelationship(ctx, *relaRtoS)
		}
		if err != nil {
			return err
		}

		return nil

	}, nil)
}
