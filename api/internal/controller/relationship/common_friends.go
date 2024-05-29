package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"log"
)

func (i impl) CommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	var user1Friends model.Users
	var user2Friends model.Users

	user1, err := i.repo.User().GetUserByEmail(ctx, email1)
	if err != nil {
		log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user1 by email %s", err, email1))
		return nil, err
	}

	user2, err := i.repo.User().GetUserByEmail(ctx, email2)
	if err != nil {
		log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user2 by email %s", err, email2))
		return nil, err
	}

	user1Friends, err = i.repo.Relationship().GetFriends(ctx, user1)
	if err != nil {
		log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user1 friends", err))
		return nil, err
	}
	user2Friends, err = i.repo.Relationship().GetFriends(ctx, user2)
	if err != nil {
		log.Printf(controller.LogErrMessage("CommonFriends", "retrieving user2 friends", err))
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
