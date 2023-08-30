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

func TestImpl_UpdateRelationship(t *testing.T) {

	testConst := struct {
		senderEmail   string
		receiverEmail string
	}{
		senderEmail:   "sender@test.com",
		receiverEmail: "receiver@test.com",
	}
	type test struct {
		statusBefore string
		statusAfter  string
		expErr       error
	}
	os.Setenv("DB_URL", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname))
	for s, tc := range map[string]test{
		"success": {
			statusBefore: orm.StatusFriends,
			statusAfter:  orm.StatusRBlockedS,
			expErr:       nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				receiverOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				}
				receiverOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

				senderOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.receiverEmail,
				}
				senderOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

				relaBeforeOrm := &orm.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     tc.statusBefore,
				}

				relaBeforeOrm.Insert(context.Background(), executor, boil.Blacklist())

				relaAfter := model.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     tc.statusAfter,
				}

				repo := New(executor)

				// When:
				err := repo.UpdateRelationship(context.Background(), relaAfter)

				// Then:
				exists, _ := relaAfter.Orm().Exists(context.Background(), executor)
				require.ErrorIs(t, tc.expErr, err)
				require.True(t, exists)
			})
		})
	}
}
