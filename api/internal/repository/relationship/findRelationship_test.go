package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
	"github.com/hxcuber/friends-management/api/pkg/testutil"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"testing"
)

func TestImpl_FindRelationship(t *testing.T) {
	testConst := struct {
		receiverEmail string
		senderEmail   string
		status        string
	}{
		"receiver@test.com",
		"sender@test.com",
		orm.StatusRSubscribedS,
	}
	type test struct {
		duplicate bool
		expErr    error
	}

	os.Setenv("DB_URL", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname))

	for s, tc := range map[string]test{
		"success": {
			duplicate: true,
			expErr:    nil,
		},
		"absence": {
			duplicate: false,
			expErr:    ErrRelationshipNotFound,
		},
	} {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				receiverOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.receiverEmail,
				}
				senderOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				}
				receiverOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
				senderOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
				if tc.duplicate {
					relaOrm := &orm.Relationship{
						ReceiverID: receiverOrm.UserID,
						SenderID:   senderOrm.UserID,
						Status:     testConst.status,
					}
					relaOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
				}
				expOut := model.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     testConst.status,
				}

				repo := New(executor)

				// When:
				out, err := repo.FindRelationship(context.Background(), model.User{
					UserID:    receiverOrm.UserID,
					UserEmail: testConst.receiverEmail,
				}, model.User{
					UserID:    senderOrm.UserID,
					UserEmail: testConst.senderEmail,
				})

				// Then:
				if tc.expErr != nil {
					require.ErrorContains(t, err, tc.expErr.Error())
				} else {
					require.ErrorIs(t, tc.expErr, err)
					require.Equal(t, expOut, *out)
				}
			})
		})
	}
}
