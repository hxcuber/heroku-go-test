package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestImpl_Receivers(t *testing.T) {
	testConst := struct {
		senderEmail string
	}{
		senderEmail: "sender@test.com",
	}
	type userRepoGetUser struct {
		out model.User
		err error
	}
	type relaRepoGetSubscribers struct {
		err error
	}
	type relaRepoGetReceivers struct {
		err error
	}

	type test struct {
		getSender       userRepoGetUser
		getSubscribers  relaRepoGetSubscribers
		getReceivers    relaRepoGetReceivers
		subscriberCount int
		receiverCount   int
		commonCount     int
		noneCount       int
		text            string
		expErr          error
	}

	for s, tc := range map[string]test{
		"no_user": {
			getSender: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 0,
			receiverCount:   0,
			commonCount:     0,
			noneCount:       0,
			text:            "",
			expErr:          user.ErrEmailNotFound,
		},
		"empty": {
			getSender: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				},
				err: nil,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 0,
			receiverCount:   0,
			commonCount:     0,
			noneCount:       0,
			text:            "hello testing testinghelloooo ",
			expErr:          nil,
		},
		"no_subscribers": {
			getSender: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				},
				err: nil,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 0,
			receiverCount:   7,
			commonCount:     9,
			noneCount:       2,
			text:            "hello testing testinghelloooo ",
			expErr:          nil,
		},
		"no_receivers": {
			getSender: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				},
				err: nil,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 5,
			receiverCount:   0,
			commonCount:     9,
			noneCount:       2,
			text:            "hello testing testinghelloooo ",
			expErr:          nil,
		},
		"no_common": {
			getSender: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				},
				err: nil,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 5,
			receiverCount:   7,
			commonCount:     0,
			noneCount:       2,
			text:            "hello testing testinghelloooo ",
			expErr:          nil,
		},
		"no_none": {
			getSender: userRepoGetUser{
				out: model.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				},
				err: nil,
			},
			getSubscribers: relaRepoGetSubscribers{
				err: nil,
			},
			getReceivers: relaRepoGetReceivers{
				err: nil,
			},
			subscriberCount: 5,
			receiverCount:   7,
			commonCount:     9,
			noneCount:       0,
			text:            "hello testing testinghelloooo ",
			expErr:          nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			userRepo := user.NewMockRepository(t)
			relaRepo := relationship.NewMockRepository(t)
			registry := repository.NewMockRegistry(t)
			relaCtrl := New(registry)

			var expOut []string
			userRepo.On("GetUserByEmail", mock.Anything, testConst.senderEmail).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.getSender.out, tc.getSender.err
				},
			)
			registry.On("User").Return(
				func() user.Repository {
					return userRepo
				})

			var stringBuilder strings.Builder
			stringBuilder.WriteString(tc.text)

			var receivers model.Users
			for i := 1; i <= tc.receiverCount; i++ {
				receiverEmail := fmt.Sprintf("receiver%d@test.com", i)
				receiver := model.User{
					UserID:    int64(i),
					UserEmail: receiverEmail,
				}

				receivers = append(receivers, receiver)
				stringBuilder.WriteString(receiverEmail + " ")
			}

			var commons model.Users
			for i := 1; i <= tc.commonCount; i++ {
				commonEmail := fmt.Sprintf("common%d@test.com", i)
				common := model.User{
					UserID:    int64(i + tc.receiverCount),
					UserEmail: commonEmail,
				}

				commons = append(commons, common)
				receivers = append(receivers, common)
				stringBuilder.WriteString(commonEmail + " ")
			}

			for i := 1; i <= tc.noneCount; i++ {
				stringBuilder.WriteString(fmt.Sprintf("none%d@test.com", i))
			}

			if tc.getSender.err == nil {
				var subscribers model.Users
				for i := 1; i <= tc.subscriberCount; i++ {
					subscriberEmail := fmt.Sprintf("subscriber%d@test.com", i)
					subscriber := model.User{
						UserID:    int64(i + tc.commonCount + tc.receiverCount),
						UserEmail: subscriberEmail,
					}

					subscribers = append(subscribers, subscriber)
					expOut = append(expOut, subscriber.UserEmail)
				}

				for _, r := range receivers {
					expOut = append(expOut, r.UserEmail)
				}
				subscribers = append(subscribers, commons...)

				relaRepo.On("GetSubscribers", mock.Anything, tc.getSender.out).Return(
					func(ctx context.Context, sender model.User) (model.Users, error) {
						return subscribers, tc.getSubscribers.err
					})

				relaRepo.On("GetReceiversFromEmails", mock.Anything, tc.getSender.out, mock.AnythingOfType("[]string")).Return(
					func(ctx context.Context, sender model.User, emailList []string) (model.Users, error) {
						return receivers, tc.getReceivers.err
					})

				registry.On("Relationship").Return(
					func() relationship.Repository {
						return relaRepo
					})
			}

			out, err := relaCtrl.Receivers(context.Background(), testConst.senderEmail, stringBuilder.String())

			require.ErrorIs(t, err, tc.expErr)
			require.ElementsMatch(t, out, expOut)
		})
	}
}
