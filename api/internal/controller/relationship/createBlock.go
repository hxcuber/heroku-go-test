package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) CreateBlock(ctx context.Context, requestorEmail string, targetEmail string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		sender, err := txRepo.Relationship().GetUserByEmail(ctx, targetEmail)
		if err != nil {
			return err
		}

		receiver, err := txRepo.Relationship().GetUserByEmail(ctx, requestorEmail)
		if err != nil {
			return err
		}

		return txRepo.Relationship().CreateBlock(ctx, sender, receiver)
	}, nil)
}
