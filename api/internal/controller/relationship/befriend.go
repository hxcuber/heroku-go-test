package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/hxcuber/friends-management/api/internal/repository/relationship"
	"github.com/pkg/errors"
	"log"
)

func (i impl) Befriend(ctx context.Context, email1 string, email2 string) error {
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.User().GetUserByEmail(ctx, email1)
		if err != nil {
			log.Printf(controller.LogErrMessage("Befriend", "retrieving user1 by email %s", err, email1))
			return err
		}

		user2, err := txRepo.User().GetUserByEmail(ctx, email2)
		if err != nil {
			log.Printf(controller.LogErrMessage("Befriend", "retrieving user2 by email %s", err, email2))
			return err
		}

		rela1to2, err := txRepo.Relationship().FindRelationship(ctx, user1, user2)
		if err != nil {
			if !errors.Is(err, relationship.ErrRelationshipNotFound) {
				log.Printf(controller.LogErrMessage("Befriend", "retrieving relationship 1 to 2", err))
				return err
			}
			rela1to2 = &model.Relationship{
				ReceiverID: user1.UserID,
				SenderID:   user2.UserID,
				Status:     orm.StatusFriends,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *rela1to2)
			if err != nil {
				log.Printf(controller.LogErrMessage("Befriend", "creating relationship 1 to 2", err))
				return err
			}
		} else {
			switch rela1to2.Status {
			case orm.StatusRBlockedS:
				{
					log.Printf(controller.LogErrMessage("Befriend", "controller logic", ErrBlocked))
					return ErrBlocked
				}
			case orm.StatusFriends:
				{
					log.Printf(controller.LogErrMessage("Befriend", "controller logic", ErrAlreadyCreated))
					return ErrAlreadyCreated
				}
			default:
				{
					rela1to2.Status = orm.StatusFriends
				}
			}
			err = txRepo.Relationship().UpdateRelationship(ctx, *rela1to2)
			if err != nil {
				log.Printf(controller.LogErrMessage("Befriend", "updating relationship 1 to 2", err))
				return err
			}
		}

		rela2to1, err := txRepo.Relationship().FindRelationship(ctx, user2, user1)
		if err != nil {
			if !errors.Is(err, relationship.ErrRelationshipNotFound) {
				log.Printf(controller.LogErrMessage("Befriend", "finding relationship 2 to 1", err))
				return err
			}
			rela2to1 = &model.Relationship{
				ReceiverID: user2.UserID,
				SenderID:   user1.UserID,
				Status:     orm.StatusFriends,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *rela2to1)
			if err != nil {
				log.Printf(controller.LogErrMessage("Befriend", "creating relationship 2 to 1", err))
				return err
			}
		} else {
			switch rela2to1.Status {
			case orm.StatusRBlockedS:
				{
					log.Printf(controller.LogErrMessage("Befriend", "controller logic", ErrBlocked))
					return ErrBlocked
				}
			case orm.StatusFriends:
				{
					log.Printf(controller.LogErrMessage("Befriend", "controller logic", ErrBlocked))
					return ErrAlreadyCreated
				}

			default:
				{
					rela2to1.Status = orm.StatusFriends
				}
			}
			err = txRepo.Relationship().UpdateRelationship(ctx, *rela2to1)
			if err != nil {
				log.Printf(controller.LogErrMessage("Befriend", "updating relationship 2 to 1", err))
				return err
			}
		}

		return nil
	}, nil)
	if err != nil {
		log.Printf(controller.LogErrMessage("Befriend", "doing in transaction", err))
		return err
	}
	return nil
}
