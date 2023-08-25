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

func TestImpl_DeleteRelationship(t *testing.T) {
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
		expErr error
	}

	os.Setenv("DB_URL", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname))

	for s, tc := range map[string]test{
		"success": {
			expErr: nil,
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

				relaOrm := &orm.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     testConst.status,
				}
				relaOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

				relaModel := model.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     testConst.status,
				}
				repo := New(executor)

				// When:
				err := repo.DeleteRelationship(context.Background(), relaModel)

				// Then:
				if tc.expErr != nil {
					require.ErrorContains(t, err, tc.expErr.Error())
				} else {
					require.ErrorIs(t, tc.expErr, err)
				}

				exists, _ := relaOrm.Exists(context.Background(), executor)
				require.False(t, false, exists)
			})
		})
	}
}
