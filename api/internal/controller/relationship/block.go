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

func (i impl) Block(ctx context.Context, requestorEmail string, targetEmail string) error {
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		receiver, err := txRepo.User().GetUserByEmail(ctx, requestorEmail)
		if err != nil {
			log.Printf(controller.LogErrMessage("Block", "retrieving receiver by email %s", err, requestorEmail))
			return err
		}

		sender, err := txRepo.User().GetUserByEmail(ctx, targetEmail)
		if err != nil {
			log.Printf(controller.LogErrMessage("Block", "retrieving sender by email %s", err, targetEmail))
			return err
		}

		relaStoR, err := txRepo.Relationship().FindRelationship(ctx, sender, receiver)
		if err != nil {
			if !errors.Is(err, relationship.ErrRelationshipNotFound) {
				log.Printf(controller.LogErrMessage("Block", "retrieving relationship sender to receiver", err))
				return err
			}
		} else {
			if relaStoR.Status == orm.StatusRBlockedS {
				log.Printf(controller.LogErrMessage("Block", "controller logic", ErrBlocked))
				return ErrBlocked
			}
			err = txRepo.Relationship().DeleteRelationship(ctx, *relaStoR)
			if err != nil {
				log.Printf(controller.LogErrMessage("Block", "deleting relationship sender to receiver", err))
				return err
			}
		}

		relaRtoS, err := txRepo.Relationship().FindRelationship(ctx, receiver, sender)
		if err != nil {
			if !errors.Is(err, relationship.ErrRelationshipNotFound) {
				log.Printf(controller.LogErrMessage("Block", "retrieving relationship receiver to sender", err))
				return err
			}
			relaRtoS = &model.Relationship{
				ReceiverID: receiver.UserID,
				SenderID:   sender.UserID,
				Status:     orm.StatusRBlockedS,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *relaRtoS)
			if err != nil {
				log.Printf(controller.LogErrMessage("Block", "creating relationship receiver to sender", err))
				return err
			}
		} else {
			if relaRtoS.Status == orm.StatusRBlockedS {
				log.Printf(controller.LogErrMessage("Block", "controller logic", ErrAlreadyCreated))
				return ErrAlreadyCreated
			}
			relaRtoS.Status = orm.StatusRBlockedS
			err = txRepo.Relationship().UpdateRelationship(ctx, *relaRtoS)
			if err != nil {
				log.Printf(controller.LogErrMessage("Block", "updating relationship receiver to sender", err))
				return err
			}
		}

		return nil
	}, nil)
	if err != nil {
		log.Printf(controller.LogErrMessage("Block", "doing in transaction", err))
		return err
	}
	return nil
}
