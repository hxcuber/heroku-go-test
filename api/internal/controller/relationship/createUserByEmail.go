package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) CreateUserByEmail(ctx context.Context, email string) error {
	return i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		return txRepo.Relationship().CreateUserByEmail(ctx, email)
	}, nil)

}
