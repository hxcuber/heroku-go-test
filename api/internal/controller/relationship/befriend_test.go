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

func TestImpl_Befriend(t *testing.T) {
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
	type relaRepoDeleteRela struct {
		err error
	}
	type test struct {
		get1       userRepoGetUser
		get2       userRepoGetUser
		find1to2   relaRepoFindRela
		find2to1   relaRepoFindRela
		create1to2 relaRepoCreateRela
		create2to1 relaRepoCreateRela
		update1to2 relaRepoUpdateRela
		update2to1 relaRepoUpdateRela
		delete1to2 relaRepoDeleteRela
		expErr     error
	}

	for s, tc := range map[string]test{
		"fail_user_1_not_found": {
			get1: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			get2:       userRepoGetUser{},
			find1to2:   relaRepoFindRela{},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     user.ErrEmailNotFound,
		},
		"fail_user_2_not_found": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{},
				err: user.ErrEmailNotFound,
			},
			find1to2:   relaRepoFindRela{},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     user.ErrEmailNotFound,
		},
		"fail_find1to2_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: errors.New("unknown"),
			},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     errors.New("unknown"),
		},
		"fail_create1to2_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{},
			create1to2: relaRepoCreateRela{
				err: errors.New("unknown"),
			},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     errors.New("unknown"),
		},
		"fail_1to2_blocked": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRBlockedS,
				},
				err: nil,
			},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     ErrBlocked,
		},
		"fail_1to2_already_created": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusFriends,
				},
				err: nil,
			},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     ErrAlreadyCreated,
		},
		"fail_update1to2_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			find2to1:   relaRepoFindRela{},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{
				err: errors.New("unknown"),
			},
			update2to1: relaRepoUpdateRela{},
			expErr:     errors.New("unknown"),
		},
		"fail_find2to1_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			find2to1: relaRepoFindRela{
				out: nil,
				err: errors.New("unknown"),
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     errors.New("unknown"),
		},
		"fail_create2to1_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			find2to1: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{
				err: errors.New("unknown"),
			},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     errors.New("unknown"),
		},
		"fail_delete2to1_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRBlockedS,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			delete1to2: relaRepoDeleteRela{
				err: errors.New("unknown"),
			},
			expErr: errors.New("unknown"),
		},
		"fail_2to1_blocked": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRBlockedS,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     ErrBlocked,
		},
		"fail_2to1_already_created": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusFriends,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     ErrAlreadyCreated,
		},
		"fail_update2to1_unknown": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{
				err: errors.New("unknown"),
			},
			expErr: errors.New("unknown"),
		},
		"success_create_update": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     nil,
		},
		"success_update_both": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 1,
					SenderID:   2,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			find2to1: relaRepoFindRela{
				out: &model.Relationship{
					ReceiverID: 2,
					SenderID:   1,
					Status:     orm.StatusRSubscribedS,
				},
				err: nil,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{
				err: nil,
			},
			update2to1: relaRepoUpdateRela{
				err: nil,
			},
			expErr: nil,
		},
		"success_create_both": {
			get1: userRepoGetUser{
				out: model.User{
					UserID:    1,
					UserEmail: testConst.email1,
				},
				err: nil,
			},
			get2: userRepoGetUser{
				out: model.User{
					UserID:    2,
					UserEmail: testConst.email2,
				},
				err: nil,
			},
			find1to2: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			find2to1: relaRepoFindRela{
				out: nil,
				err: relationship.ErrRelationshipNotFound,
			},
			create1to2: relaRepoCreateRela{},
			create2to1: relaRepoCreateRela{},
			update1to2: relaRepoUpdateRela{},
			update2to1: relaRepoUpdateRela{},
			expErr:     nil,
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

			userRepo.On("GetUserByEmail", mock.Anything, testConst.email1).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.get1.out, tc.get1.err
				},
			)
			if tc.get1.err != nil {
				goto execute
			}

			userRepo.On("GetUserByEmail", mock.Anything, testConst.email2).Return(
				func(ctx context.Context, email string) (model.User, error) {
					return tc.get2.out, tc.get2.err
				},
			)
			if tc.get2.err != nil {
				goto execute
			}

			registry.On("Relationship").Return(
				func() relationship.Repository {
					return relaRepo
				})

			relaRepo.On("FindRelationship", mock.Anything, tc.get1.out, tc.get2.out).Return(
				func(ctx context.Context, user1 model.User, user2 model.User) (*model.Relationship, error) {
					return tc.find1to2.out, tc.find1to2.err
				})
			if tc.find1to2.err != nil {
				if !errors.Is(tc.find1to2.err, relationship.ErrRelationshipNotFound) {
					goto execute
				}
				relaRepo.On("CreateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.get1.out.UserID,
					SenderID:   tc.get2.out.UserID,
					Status:     orm.StatusFriends,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.create1to2.err
					})
				if tc.create1to2.err != nil {
					goto execute
				}
			} else {
				if tc.find1to2.out.Status != orm.StatusRSubscribedS {
					goto execute
				}
				relaRepo.On("UpdateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.find1to2.out.ReceiverID,
					SenderID:   tc.find1to2.out.SenderID,
					Status:     orm.StatusFriends,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.update1to2.err
					})
				if tc.update1to2.err != nil {
					goto execute
				}
			}

			relaRepo.On("FindRelationship", mock.Anything, tc.get2.out, tc.get1.out).Return(
				func(ctx context.Context, user2 model.User, user1 model.User) (*model.Relationship, error) {
					return tc.find2to1.out, tc.find2to1.err
				})
			if tc.find2to1.err != nil {
				if !errors.Is(tc.find2to1.err, relationship.ErrRelationshipNotFound) {
					goto execute
				}
				relaRepo.On("CreateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.get2.out.UserID,
					SenderID:   tc.get1.out.UserID,
					Status:     orm.StatusFriends,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.create2to1.err
					})
				if tc.create2to1.err != nil {
					goto execute
				}
			} else {
				if tc.find2to1.out.Status == orm.StatusRBlockedS {
					relaRepo.On("DeleteRelationship", mock.Anything, model.Relationship{
						ReceiverID: tc.get1.out.UserID,
						SenderID:   tc.get2.out.UserID,
						Status:     orm.StatusFriends,
					}).Return(
						func(ctx context.Context, relationship model.Relationship) error {
							return tc.delete1to2.err
						})
					goto execute
				}
				if tc.find2to1.out.Status == orm.StatusFriends {
					goto execute
				}
				relaRepo.On("UpdateRelationship", mock.Anything, model.Relationship{
					ReceiverID: tc.find2to1.out.ReceiverID,
					SenderID:   tc.find2to1.out.SenderID,
					Status:     orm.StatusFriends,
				}).Return(
					func(ctx context.Context, relationship model.Relationship) error {
						return tc.update2to1.err
					})
				if tc.update2to1.err != nil {
					goto execute
				}
			}

		execute:
			err := relaCtrl.Befriend(context.Background(), testConst.email1, testConst.email2)
			if tc.expErr != nil {
				require.ErrorContains(t, err, tc.expErr.Error())
			} else {
				require.ErrorIs(t, err, tc.expErr)
			}
		})
	}
}
