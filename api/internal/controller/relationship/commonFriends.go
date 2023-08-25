package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"log"
)

func (i impl) CommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	var user1Friends model.Users
	var user2Friends model.Users
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		user1, err := txRepo.Relationship().GetUserByEmail(ctx, email1)
		if err != nil {
			log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user1 by email %s", err, email1))
			return err
		}

		user2, err := txRepo.Relationship().GetUserByEmail(ctx, email2)
		if err != nil {
			log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user2 by email %s", err, email2))
			return err
		}

		user1Friends, err = txRepo.Relationship().GetFriends(ctx, user1)
		if err != nil {
			log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user1 friends", err))
			return err
		}
		user2Friends, err = txRepo.Relationship().GetFriends(ctx, user2)
		if err != nil {
			log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user2 friends", err))
			return err
		}

		return nil
	}, nil)
	if err != nil {
		log.Printf(controller.LogErrMessage("CommonFriends", "doing in transaction", err))
		return nil, err
	}

	// This goes against convention but is beneficial for the JSON stage.
	commonFriendsEmail := []string{}

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
