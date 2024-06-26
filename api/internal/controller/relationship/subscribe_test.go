package relationship

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/hxcuber/friends-management/api/internal/repository/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImpl_Subscribe(t *testing.T) {
	testConst := struct {
		requestorEmail string
		targetEmail    string
	}{
		requestorEmail: "request@test.com",
		targetEmail:    "target@test.com",
	}
	type userRepoGetUser struct {
		out model.User
		err error
	}
	type relaRepoFindRela struct {
		out *model.Relationship
		err error
	}
	type relaRepoCreateRela struct {
		err error
	}

	type relaRepoUpdateRela struct {
		err error
	}
	type test struct {
		getReq     userRepoGetUser
		getTar     userRepoGetUser
		findStoR   relaRepoFindRela
		findRtoS   relaRepoFindRela
		createRtoS relaRepoCreateRela
		updateRtoS relaRepoUpdateRela
		expErr     error
	}

	for s, tc := range map[string]test{
		"success_update": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.targetEmail,
				},
				err: nil,
			},
			findStoR: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			findRtoS: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRBlockedS,
				},
				err: nil,
			},
			createRtoS: relaRepoCreateRela{nil},
			updateRtoS: relaRepoUpdateRela{nil},
			expErr:     nil,
		},
		"success_create": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.targetEmail,
				},
				err: nil,
			},
			findStoR: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			findRtoS: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			createRtoS: relaRepoCreateRela{nil},
			updateRtoS: relaRepoUpdateRela{nil},
			expErr:     nil,
		},
		"fail_req_not_found": {
			getReq: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			getTar:     userRepoGetUser{},
			findStoR:   relaRepoFindRela{},
			findRtoS:   relaRepoFindRela{},
			createRtoS: relaRepoCreateRela{},
			updateRtoS: relaRepoUpdateRela{},
			expErr:     user.ErrEmailNotFound,
		},
		"fail_tar_not_found": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			findStoR:   relaRepoFindRela{},
			findRtoS:   relaRepoFindRela{},
			createRtoS: relaRepoCreateRela{},
			updateRtoS: relaRepoUpdateRela{},
			expErr:     user.ErrEmailNotFound,
		},
		"fail_already_created": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.targetEmail,
				},
				err: nil,
			},
			findStoR: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			findRtoS: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			createRtoS: relaRepoCreateRela{nil},
			updateRtoS: relaRepoUpdateRela{nil},
			expErr:     ErrAlreadyCreated,
		},
		"fail_already_friends": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.targetEmail,
				},
				err: nil,
			},
			findStoR: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 2,
					SenderID:   1,
					Status:     orm.StatusFriends,
				},
				err: nil,
			},
			findRtoS: relaRepoFindRela{
				out: nil,
				err: nil,
			},
			createRtoS: relaRepoCreateRela{nil},
			updateRtoS: relaRepoUpdateRela{nil},
			expErr:     ErrFriends,
		},
		"fail_blocked": {
			getReq: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.requestorEmail,
				},
				err: nil,
			},
			getTar: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.targetEmail,
				},
				err: nil,
			},
			findStoR: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 2,
					SenderID:   1,
					Status:     orm.StatusRBlockedS,
				},
				err: nil,
			},
			findRtoS: relaRepoFindRela{
				out: nil,
				err: nil,
			},
			createRtoS: relaRepoCreateRela{nil},
			updateRtoS: relaRepoUpdateRela{nil},
			expErr:     ErrBlocked,
		},
	} {
		t.Run(s, func(t *testing.T) {
			userRepo := user.NewMockRepository(t)
			relaRepo := relationship.NewMockRepository(t)
			registry := repository.NewMockRegistry(t)
			relaCtrl := New(registry)

			registry.On("DoInTx", mock.Anything, mock.Anything, nil).Return(
				func(ctx context.Context, txFunc func(ctx context.Context, txRepo repository.Registry) error, policy backoff.BackOff) error {
					return txFunc(ctx, registry)
				})
			registry.On("User").Return(
				func() user.Repository {
					return userRepo
				})

			userRepo.On("GetUserByEmail", mock.Anything, testConst.requestorEmail).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.getReq.out, tc.getReq.err
				},
			)
			if tc.getReq.err != nil {
				goto execute
			}

			userRepo.On("GetUserByEmail", mock.Anything, testConst.targetEmail).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.getTar.out, tc.getTar.err
				},
			)
			if tc.getTar.err != nil {
				goto execute
			}

			registry.On("Relationship").Return(
				func() relationship.Repository {
					return relaRepo
				})

			relaRepo.On("FindRelationship", mock.Anything, tc.getTar.out, tc.getReq.out).Return(
				func(ctx context.Context, receiver model.User, sender model.User) (*model.Relationship, error) {
					return tc.findStoR.out, tc.findStoR.err
				})
			if tc.findStoR.err != nil {
				if !errors.Is(tc.findStoR.err, relationship.ErrRelationshipNotFound) {
					goto execute
				}
			} else if tc.findStoR.out.Status == orm.StatusRBlockedS {
				goto execute
			} else if tc.findStoR.out.Status == orm.StatusFriends {
				goto execute
			}

			relaRepo.On("FindRelationship", mock.Anything, tc.getReq.out, tc.getTar.out).Return(
				func(ctx context.Context, sender model.User, receiver model.User) (*model.Relationship, error) {
					return tc.findRtoS.out, tc.findRtoS.err
				})
			if tc.findRtoS.err != nil {
				if !errors.Is(tc.findRtoS.err, relationship.ErrRelationshipNotFound) {
					goto execute
				}
				relaRepo.On("CreateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.getReq.out.UserID,
					SenderID:   tc.getTar.out.UserID,
					Status:     orm.StatusRSubscribedS,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.createRtoS.err
					})
			} else {
				if tc.findRtoS.out.Status != orm.StatusRBlockedS {
					goto execute
				}
				relaRepo.On("UpdateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.findRtoS.out.ReceiverID,
					SenderID:   tc.findRtoS.out.SenderID,
					Status:     orm.StatusRSubscribedS,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.updateRtoS.err
					})
			}

		execute:
			err := relaCtrl.Subscribe(context.Background(), testConst.requestorEmail, testConst.targetEmail)
			require.ErrorIs(t, err, tc.expErr)
		})
	}
}
