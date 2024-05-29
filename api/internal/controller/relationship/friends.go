package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"log"
)

func (i impl) Friends(ctx context.Context, email string) ([]string, error) {

	var friends model.Users

	user, err := i.repo.User().GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf(controller.LogErrMessage("Friends", "retrieving user by email %s", err, email))
		return nil, err
	}

	friends, err = i.repo.Relationship().GetFriends(ctx, user)
	if err != nil {
		log.Printf(controller.LogErrMessage("Friends", "retrieving user friends", err))
		return nil, err
	}

	if err != nil {
		log.Printf(controller.LogErrMessage("Friends", "doing in transaction", err))
		return nil, err
	}

	// This goes against convention but is beneficial for the JSON stage.
	friendsEmail := []string{}

	for _, friend := range friends {
		friendsEmail = append(friendsEmail, friend.UserEmail)
	}
	return friendsEmail, nil
}
