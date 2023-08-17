package relationship

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/pkg/errors"
)

func (i impl) CreateConnection(ctx context.Context, email1 string, email2 string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			return err
		}

		relationship, err := txRepo.Relationship().FindRelationship(ctx, user1, user2)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return txRepo.Relationship().CreateConnection(ctx, user1, user2)
		}

		if relationship.Friends == true {
			return controller.ErrAlreadyCreated
		}

		return txRepo.Relationship().CreateConnection(ctx, user1, user2)
	}, nil)
}
