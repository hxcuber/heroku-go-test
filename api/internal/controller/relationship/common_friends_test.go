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

func TestImpl_CommonFriends(t *testing.T) {
	testConst := struct {
		email1 string
		email2 string
	}{
		email1: "test1@test.com",
		email2: "test2@test.com",
	}
	type userRepoGetUser struct {
		out model.User
		err error
	}
	type relaRepoGetFriends struct {
		err error
	}
	type test struct {
		getUser1     userRepoGetUser
		getUser2     userRepoGetUser
		getFriends1  relaRepoGetFriends
		getFriends2  relaRepoGetFriends
		friend1Count int
		friend2Count int
		commonCount  int
		expErr       error
	}

	for s, tc := range map[string]test{
		"no_user_1": {
			getUser1: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			getUser2: userRepoGetUser{
				out: model.User{},
				err: nil,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 0,
			friend2Count: 0,
			commonCount:  0,
			expErr:       user.ErrEmailNotFound,
		},
		"no_user_2": {
			getUser1: userRepoGetUser{
				out: model.User{},
				err: nil,
			},
			getUser2: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 0,
			friend2Count: 0,
			commonCount:  0,
			expErr:       user.ErrEmailNotFound,
		},
		"0_common": {
			getUser1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			getUser2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 7,
			friend2Count: 5,
			commonCount:  0,
			expErr:       nil,
		},
		"1_common": {
			getUser1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			getUser2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 7,
			friend2Count: 9,
			commonCount:  1,
			expErr:       nil,
		},
		"2_friends": {
			getUser1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			getUser2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 7,
			friend2Count: 5,
			commonCount:  2,
			expErr:       nil,
		},
		"10_friends": {
			getUser1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			getUser2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			getFriends1: relaRepoGetFriends{
				err: nil,
			},
			getFriends2: relaRepoGetFriends{
				err: nil,
			},
			friend1Count: 7,
			friend2Count: 5,
			commonCount:  2,
			expErr:       nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			userRepo := user.NewMockRepository(t)
			relaRepo := relationship.NewMockRepository(t)
			registry := repository.NewMockRegistry(t)
			var expOut []string
			userRepo.On("GetUserByEmail", mock.Anything, testConst.email1).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.getUser1.out, tc.getUser1.err
				},
			)

			if tc.getUser1.err == nil {
				userRepo.On("GetUserByEmail", mock.Anything, testConst.email2).Return(
					func(ctx context.Context, email string) (model.User, error) {
						return tc.getUser2.out, tc.getUser2.err
					},
				)
			}

			registry.On("User").Return(
				func() user.Repository {
					return userRepo
				})

			if tc.getUser1.err == nil && tc.getUser2.err == nil {
				var friends1 model.Users
				for i := 1; i <= tc.friend1Count; i++ {
					friendEmail := fmt.Sprintf("1friend%d@test.com", i)
					friend := model.User{
						UserID:    int64(i + 2),
						UserEmail: friendEmail,
					}

					friends1 = append(friends1, friend)
				}

				var friends2 model.Users
				for i := 1; i <= tc.friend2Count; i++ {
					friendEmail := fmt.Sprintf("2friend%d@test.com", i)
					friend := model.User{
						UserID:    int64(i + tc.friend1Count + 2),
						UserEmail: friendEmail,
					}

					friends2 = append(friends2, friend)
				}

				for i := 1; i <= tc.commonCount; i++ {
					commonEmail := fmt.Sprintf("common%d@test.com", i)
					common := model.User{
						UserID:    int64(i + tc.friend1Count + tc.friend2Count + 2),
						UserEmail: commonEmail,
					}

					friends1 = append(friends1, common)
					friends2 = append(friends2, common)
					expOut = append(expOut, commonEmail)
				}

				relaRepo.On("GetFriends", mock.Anything, tc.getUser1.out).Return(
					func(ctx context.Context, user model.User) (model.Users, error) {
						return friends1, tc.getFriends1.err
					})
				relaRepo.On("GetFriends", mock.Anything, tc.getUser2.out).Return(
					func(ctx context.Context, user model.User) (model.Users, error) {
						return friends2, tc.getFriends1.err
					})
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
			out, err := relaCtrl.CommonFriends(context.Background(), testConst.email1, testConst.email2)

			require.ErrorIs(t, err, tc.expErr)
			require.ElementsMatch(t, out, expOut)
		})
	}
}
