package relationship

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
)

func (i impl) PostCreateSubscription(ctx context.Context, requestorEmail string, targetEmail string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		sender, err := txRepo.Relationship().GetUserByEmail(ctx, targetEmail)
		if err != nil {
			return err
		}

		receiver, err := txRepo.Relationship().GetUserByEmail(ctx, requestorEmail)
		if err != nil {
			return err
		}

		relationship, err := txRepo.Relationship().FindRelationship(ctx, sender, receiver)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return txRepo.Relationship().CreateSubscription(ctx, sender, receiver)
		}

		if relationship.Status == orm.SubscriptionStatusRSubscribedS {
			return controller.ErrAlreadyCreated
		}

		return txRepo.Relationship().CreateSubscription(ctx, sender, receiver)
	}, nil)
}
