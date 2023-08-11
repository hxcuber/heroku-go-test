package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) CreateUserByEmail(ctx context.Context, email string) error {
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		err := txRepo.Relationship().CreateUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		return nil
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
