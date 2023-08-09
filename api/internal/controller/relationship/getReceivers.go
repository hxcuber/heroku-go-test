package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/pkg/util"
	"strings"
)

func (i impl) GetReceivers(ctx context.Context, senderEmail string, text string) ([]string, error) {
	tokens := strings.Split(strings.ReplaceAll(text, ",", " "), " ")

	var emailsInText []string
	for _, token := range tokens {
		if util.IsEmail(token) {
			emailsInText = append(emailsInText, token)
		}
	}

	var validEmailsMentioned []string

	var subscribers model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		var err error
		subscribers, err = i.repo.Relationship().GetSubscribers(ctx, senderEmail)
		if err != nil {
			return err
		}

		for _, email := range emailsInText {
			status, err := i.repo.Relationship().GetNotificationStatus(ctx, senderEmail, email)
			if err == nil || status != "r_blocked_s" {
				validEmailsMentioned = append(validEmailsMentioned, email)
			}
		}

		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	var receiversEmail []string

	for _, subscriber := range subscribers {
		receiversEmail = append(receiversEmail, subscriber.UserEmail)
	}

	for _, email := range validEmailsMentioned {
		receiversEmail = append(receiversEmail, email)
	}

	return receiversEmail, nil
}
