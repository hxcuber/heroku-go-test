package user

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"log"
)

func (i impl) CreateUserByEmail(ctx context.Context, email string) error {
	err := i.repo.DoInTx(context.Background(), func(ctx context.Context, txRepo repository.Registry) error {
		_, err := txRepo.User().GetUserByEmail(ctx, email)
		if err != nil {
			if !errors.Is(err, user.ErrEmailNotFound) {
				log.Printf(controller.LogErrMessage("CreateUserByEmail", "retrieving user by email", err))
				return err
			}
			err = txRepo.User().CreateUserByEmail(ctx, email)
			if err != nil {
				log.Printf(controller.LogErrMessage("CreateUserByEmail", "creating user by email", err))
				return err
			}
			return nil
		}
		log.Printf(controller.LogErrMessage("CreateUserByEmail", "controller logic", ErrAlreadyCreated))
		return ErrAlreadyCreated
	}, nil)
	if err != nil {
		log.Printf(controller.LogErrMessage("CreateUserByEmail", "doing in transaction", err))
		return err
	}
	return nil
}
