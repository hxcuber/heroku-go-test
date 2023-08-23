package relationship

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
	"log"
)

func (i impl) Subscribe(ctx context.Context, requestorEmail string, targetEmail string) error {
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		receiver, err := txRepo.Relationship().GetUserByEmail(ctx, requestorEmail)
		if err != nil {
			log.Printf(LogErrMessage("Subscribe", "retrieving receiver by email %s", err, requestorEmail))
			return err
		}

		sender, err := txRepo.Relationship().GetUserByEmail(ctx, targetEmail)
		if err != nil {
			log.Printf(LogErrMessage("Subscribe", "retrieving sender by email %s", err, targetEmail))
			return err
		}

		relaStoR, err := txRepo.Relationship().FindRelationship(ctx, sender, receiver)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf(LogErrMessage("Subscribe", "retrieving relationship sender to receiver", err))
				return err
			}
		} else {
			if relaStoR.Status == orm.StatusRBlockedS {
				log.Printf(LogErrMessage("Subscribe", "controller logic", ErrBlocked))
				return ErrBlocked
			}
		}

		relaRtoS, err := txRepo.Relationship().FindRelationship(ctx, receiver, sender)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf(LogErrMessage("Subscribe", "retrieving relationship receiver to sender", err))
				return err
			}
			relaRtoS = &model.Relationship{
				ReceiverID: receiver.UserID,
				SenderID:   sender.UserID,
				Status:     orm.StatusRSubscribedS,
			}
			err = txRepo.Relationship().CreateRelationship(ctx, *relaRtoS)
			if err != nil {
				log.Printf(LogErrMessage("Subscribe", "creating relationship receiver to sender", err))
				return err
			}
		} else {
			switch relaRtoS.Status {
			case orm.StatusRSubscribedS:
				{
					log.Printf(LogErrMessage("Subscribe", "controller logic", ErrAlreadyCreated))
					return ErrAlreadyCreated
				}
			case orm.StatusFriends:
				{
					log.Printf(LogErrMessage("Subscribe", "controller logic", ErrFriends))
					return ErrFriends
				}
			default:
				{
					relaRtoS.Status = orm.StatusRSubscribedS
				}
			}
			err = txRepo.Relationship().UpdateRelationship(ctx, *relaRtoS)
			if err != nil {
				log.Printf(LogErrMessage("Subscribe", "updating relationship receiver to sender", err))
				return err
			}
		}

		return nil
	}, nil)
	if err != nil {
		log.Printf(LogErrMessage("Subscribe", "doing in transaction", err))
		return err
	}
	return nil
}
