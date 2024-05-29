package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/pkg/util"
	"log"
	"strings"
)

func (i impl) Receivers(ctx context.Context, senderEmail string, text string) ([]string, error) {
	// Replacing "'" with " " prevents the user from escaping the "'" used in the queries
	// preventing SQL injection.
	tokens := strings.Split(strings.ReplaceAll(strings.ReplaceAll(text, ",", " "), "'", " "), " ")
	var emailList []string
	for _, token := range tokens {
		if util.IsEmail(token) {
			emailList = append(emailList, token)
		}
	}

	var validUsersMentioned model.Users
	var subscribers model.Users

	sender, err := i.repo.User().GetUserByEmail(ctx, senderEmail)
	if err != nil {
		log.Printf(controller.LogErrMessage("Receivers", "retrieving sender by email %s", err, senderEmail))
		return nil, err
	}

	subscribers, err = i.repo.Relationship().GetSubscribers(ctx, sender)
	if err != nil {
		log.Printf(controller.LogErrMessage("Receivers", "retrieving sender subscribers", err))
		return nil, err
	}

	validUsersMentioned, err = i.repo.Relationship().GetReceiversFromEmails(ctx, sender, emailList)
	if err != nil {
		log.Printf(controller.LogErrMessage("Receivers", "retrieving sender receivers from email list", err))
		return nil, err
	}
	if err != nil {
		log.Printf(controller.LogErrMessage("Receivers", "doing in transaction", err))
		return nil, err
	}

	// This goes against convention but is beneficial for the JSON stage.
	receiversEmail := []string{}

	// Removing duplicates
	hash := make(map[string]bool)

	for _, subscriber := range subscribers {
		hash[subscriber.UserEmail] = true
		receiversEmail = append(receiversEmail, subscriber.UserEmail)
	}

	for _, validUserMentioned := range validUsersMentioned {
		if !hash[validUserMentioned.UserEmail] {
			receiversEmail = append(receiversEmail, validUserMentioned.UserEmail)
			hash[validUserMentioned.UserEmail] = true
		}
	}

	return receiversEmail, nil
}
