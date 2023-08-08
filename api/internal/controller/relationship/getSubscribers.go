package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
)

func (i impl) GetSubscribers(ctx context.Context, email string) ([]string, error) {
	var subscribers model.UserSlice
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		var err error
		subscribers, err = i.repo.Relationship().GetSubscribers(ctx, email)
		return err
	}, nil)
	if err != nil {
		return nil, err
	}
	var subscribersEmail []string

	for _, subscriber := range subscribers {
		subscribersEmail = append(subscribersEmail, subscriber.UserEmail)
	}
	return subscribersEmail, nil
}
