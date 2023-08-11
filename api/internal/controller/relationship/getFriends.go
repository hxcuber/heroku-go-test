package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) GetFriends(ctx context.Context, email string) ([]string, error) {
	var friends model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		var err error
		friends, err = txRepo.Relationship().GetFriends(ctx, email)
		return err
	}, nil)
	if err != nil {
		return nil, err
	}
	var friendsEmail []string

	for _, friend := range friends {
		friendsEmail = append(friendsEmail, friend.UserEmail)
	}
	return friendsEmail, nil
}
