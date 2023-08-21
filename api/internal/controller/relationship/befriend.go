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

func (i impl) Befriend(ctx context.Context, email1 string, email2 string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			return err
		}

		rela1to2, err := txRepo.Relationship().FindRelationship(ctx, user1, user2)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			rela1to2 = &model.Relationship{
				ReceiverID: user1.UserID,
				SenderID:   user2.UserID,
				Status:     orm.StatusFriends,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *rela1to2)
		} else {
			switch rela1to2.Status {
			case orm.StatusRBlockedS:
				{
					return controller.ErrBlocked
				}
			case orm.StatusFriends:
				{
					return controller.ErrAlreadyCreated
				}
			default:
				{
					rela1to2.Status = orm.StatusFriends
				}
			}
			_, err = txRepo.Relationship().UpdateRelationship(ctx, *rela1to2)
		}
		if err != nil {
			return err
		}

		rela2to1, err := txRepo.Relationship().FindRelationship(ctx, user2, user1)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			rela2to1 = &model.Relationship{
				ReceiverID: user2.UserID,
				SenderID:   user1.UserID,
				Status:     orm.StatusFriends,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *rela2to1)
		} else {
			switch rela2to1.Status {
			case orm.StatusRBlockedS:
				{
					return controller.ErrBlocked
				}
			case orm.StatusFriends:
				{
					return controller.ErrAlreadyCreated
				}
			default:
				{
					rela2to1.Status = orm.StatusFriends
				}
			}
			_, err = txRepo.Relationship().UpdateRelationship(ctx, *rela2to1)
		}
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
