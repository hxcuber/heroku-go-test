package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) GetCommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	var user1Friends model.UserSlice
	var user2Friends model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			return err
		}

		user1Friends, err = txRepo.Relationship().GetFriends(ctx, user1)
		user2Friends, err = txRepo.Relationship().GetFriends(ctx, user2)
		return err
	}, nil)
	if err != nil {
		return nil, err
	}
	var commonFriendsEmail []string

	hash := make(map[string]bool)
	for _, friend := range user1Friends {
		hash[friend.UserEmail] = true
	}

	for _, friend := range user2Friends {
		if hash[friend.UserEmail] {
			commonFriendsEmail = append(commonFriendsEmail, friend.UserEmail)
		}
	}

	return commonFriendsEmail, nil
}
