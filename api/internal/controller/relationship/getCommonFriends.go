package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) GetCommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	var commonFriends model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			return err
		}

		commonFriends, err = txRepo.Relationship().GetCommonFriends(ctx, user1, user2)
		return err
	}, nil)
	if err != nil {
		return nil, err
	}
	var commonFriendsEmail []string

	for _, friend := range commonFriends {
		commonFriendsEmail = append(commonFriendsEmail, friend.UserEmail)
	}

	return commonFriendsEmail, nil
}
