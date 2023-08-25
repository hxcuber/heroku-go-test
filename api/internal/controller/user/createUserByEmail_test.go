package user

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImpl_CreateUserByEmail(t *testing.T) {
	testConst := struct {
		email string
	}{
		email: "test@test.com",
	}
	type test struct {
		repoCreateUserByEmail error
		repoGetUserByEmail    error
		doInTxErr             error
		expErr                error
	}

	for s, tc := range map[string]test{
		"success": {
			repoCreateUserByEmail: nil,
			repoGetUserByEmail:    user.ErrEmailNotFound,
			doInTxErr:             nil,
			expErr:                nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			userRepo := user.MockRepository{}
			userRepo.On("CreateUserByEmail",
				mock.Anything, testConst.email).Return(tc.repoCreateUserByEmail)
			//func(ctx context.Context, email string) error {
			//	return tc.repoCreateUserByEmail
			//})
			userRepo.On("GetUserByEmail",
				mock.Anything, testConst.email).Return(model.User{}, tc.repoGetUserByEmail)
			//func(ctx context.Context, email string) (model.User, error) {
			//	return model.User{}, tc.repoGetUserByEmail
			//})
			reg := repository.MockRegistry{}
			reg.On("User").Return(&userRepo)
			//	func() user.Repository {
			//	return &userRepo
			//})
			reg.On("DoInTx",
				mock.Anything,
				mock.Anything,
				nil).Return( //tc.doInTxErr)
				func(ctx context.Context, f func(ctx2 context.Context, txRepo repository.Registry) error, policy backoff.BackOff) error {
					return f(ctx, &reg)
				})

			userCtrl := New(&reg)
			err := userCtrl.CreateUserByEmail(context.Background(), testConst.email)
			require.ErrorIs(t, tc.expErr, err)
		})
	}
}
