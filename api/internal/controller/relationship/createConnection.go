package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) CreateConnection(ctx context.Context, email1 string, email2 string) (bool, error) {
	var created bool
	return created, i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			return err
		}

		created, err = txRepo.Relationship().CreateConnection(ctx, user1, user2)

		return err
	}, nil)
}
