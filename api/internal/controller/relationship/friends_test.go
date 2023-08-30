package relationship

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImpl_Friends(t *testing.T) {
	testConst := struct {
		email string
	}{
		email: "test@test.com",
	}
	type userRepoGetUser struct {
		out model.User
		err error
	}
	type relaRepoGetFriends struct {
		err error
	}
	type test struct {
		getUser     userRepoGetUser
		getFriends  relaRepoGetFriends
		friendCount int
		expErr      error
	}

	for s, tc := range map[string]test{
		"no_user": {
			getUser: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			getFriends: relaRepoGetFriends{
				err: nil,
			},
			friendCount: 0,
			expErr:      user.ErrEmailNotFound,
		},
		"0_friends": {
			getUser: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.email,
				},
				err: nil,
			},
			getFriends: relaRepoGetFriends{
				err: nil,
			},
			friendCount: 0,
			expErr:      nil,
		},
		"1_friend": {
			getUser: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.email,
				},
				err: nil,
			},
			getFriends: relaRepoGetFriends{
				err: nil,
			},
			friendCount: 1,
			expErr:      nil,
		},
		"2_friends": {
			getUser: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.email,
				},
				err: nil,
			},
			getFriends: relaRepoGetFriends{
				err: nil,
			},
			friendCount: 2,
			expErr:      nil,
		},
		"10_friends": {
			getUser: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.email,
				},
				err: nil,
			},
			getFriends: relaRepoGetFriends{
				err: nil,
			},
			friendCount: 10,
			expErr:      nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			userRepo := user.NewMockRepository(t)
			relaRepo := relationship.NewMockRepository(t)
			registry := repository.NewMockRegistry(t)
			var expOut []string
			userRepo.On("GetUserByEmail", mock.Anything, testConst.email).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.getUser.out, tc.getUser.err
				},
			)
			registry.On("User").Return(
				func() user.Repository {
					return userRepo
				})

			if tc.getUser.err == nil {
				var friends model.Users
				for i := 1; i <= tc.friendCount; i++ {
					friendEmail := fmt.Sprintf("friend%d@test.com", i)
					friend := model.User{
						UserID:    int64(i),
						UserEmail: friendEmail,
					}

					friends = append(friends, friend)
					expOut = append(expOut, friend.UserEmail)
				}
				relaRepo.On("GetFriends", mock.Anything, tc.getUser.out).Return(
					func(ctx context.Context, user model.User) (model.Users, error) {
						return friends, tc.getFriends.err
					},
				)
				registry.On("Relationship").Return(
					func() relationship.Repository {
						return relaRepo
					})
			}

			registry.On("DoInTx", mock.Anything, mock.Anything, nil).Return(
				func(ctx context.Context, txFunc func(ctx context.Context, txRepo repository.Registry) error, policy backoff.BackOff) error {
					return txFunc(ctx, registry)
				})

			relaCtrl := New(registry)
			out, err := relaCtrl.Friends(context.Background(), testConst.email)

			require.ErrorIs(t, err, tc.expErr)
			require.ElementsMatch(t, out, expOut)
		})
	}
}
