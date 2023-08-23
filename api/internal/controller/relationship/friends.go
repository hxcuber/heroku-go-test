package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"log"
)

func (i impl) Friends(ctx context.Context, email string) ([]string, error) {

	var friends model.Users
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user, err := txRepo.Relationship().GetUserByEmail(ctx, email)
		if err != nil {
			log.Printf(LogErrMessage("Friends", "retrieving user by email %s", err, email))
			return err
		}

		friends, err = txRepo.Relationship().GetFriends(ctx, user)
		if err != nil {
			log.Printf(LogErrMessage("Friends", "retrieving user friends", err))
			return err
		}

		return nil
	}, nil)
	if err != nil {
		log.Printf(LogErrMessage("Friends", "doing in transaction", err))
		return nil, err
	}

	// This goes against convention but is beneficial for the JSON stage.
	friendsEmail := []string{}

	for _, friend := range friends {
		friendsEmail = append(friendsEmail, friend.UserEmail)
	}
	return friendsEmail, nil
}
