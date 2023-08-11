package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/pkg/util"
	"strings"
)

func (i impl) GetReceivers(ctx context.Context, senderEmail string, text string) ([]string, error) {
	// Replacing "'" with " " prevents the user from escaping the "'" used in the queries
	// preventing SQL injection.
	tokens := strings.Split(strings.ReplaceAll(strings.ReplaceAll(text, ",", " "), "'", " "), " ")
	var emailList []string
	for _, token := range tokens {
		if util.IsEmail(token) {
			emailList = append(emailList, token)
		}
	}

	var validUsersMentioned model.UserSlice
	var subscribers model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		var err error
		subscribers, err = txRepo.Relationship().GetSubscribers(ctx, senderEmail)
		if err != nil {
			return err
		}

		validUsersMentioned, err = txRepo.Relationship().GetReceiversFromEmails(ctx, senderEmail, emailList)
		if err != nil {
			return err
		}
		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	var receiversEmail []string

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
