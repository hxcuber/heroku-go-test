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
	type userRepoGetUser struct {
		err error
	}
	type userRepoCreateUser struct {
		err error
	}
	type test struct {
		getUser    userRepoGetUser
		createUser userRepoCreateUser
		expErr     error
	}

	for s, tc := range map[string]test{
		"success": {
			getUser: userRepoGetUser{
				err: user.ErrEmailNotFound,
			},
			createUser: userRepoCreateUser{
				err: nil,
			},
			expErr: nil,
		},
		"already_created": {
			getUser: userRepoGetUser{
				err: nil,
			},
			createUser: userRepoCreateUser{
				err: ErrAlreadyCreated,
			},
			expErr: ErrAlreadyCreated,
		},
	} {
		t.Run(s, func(t *testing.T) {
			//userRepo := user.MockRepository{}
			userRepo := user.NewMockRepository(t)
			userRepo.On("GetUserByEmail",
				mock.Anything, testConst.email).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return model.User{}, tc.getUser.err
				})
			if tc.getUser.err != nil {
				userRepo.On("CreateUserByEmail",
					mock.Anything, testConst.email).Return(
					func(ctx context.Context, email string) error {
						return tc.createUser.err
					})
			}
			registry := repository.NewMockRegistry(t)
			registry.On("User").Return(
				func() user.Repository {
					return userRepo
				})
			registry.On("DoInTx", mock.Anything, mock.Anything, nil).Return(
				func(ctx context.Context, txFunc func(ctx context.Context, txRepo repository.Registry) error, policy backoff.BackOff) error {
					return txFunc(ctx, registry)
				})

			userCtrl := New(registry)
			err := userCtrl.CreateUserByEmail(context.Background(), testConst.email)
			require.ErrorIs(t, err, tc.expErr)
		})
	}
}
