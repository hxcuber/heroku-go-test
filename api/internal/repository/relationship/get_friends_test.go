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

func TestImpl_GetFriends(t *testing.T) {
	testConst := struct {
		email string
	}{
		email: "test@test.com",
	}
	type test struct {
		friendCount int
		expErr      error
	}

	os.Setenv("DB_URL", "postgres://friends_management:@pg:5432/friends_management?sslmode=disable")
	for s, tc := range map[string]test{
		"0_friends": {
			friendCount: 0,
			expErr:      nil,
		},
		"1_friend": {
			friendCount: 1,
			expErr:      nil,
		},
		"2_friends": {
			friendCount: 2,
			expErr:      nil,
		},
		"10_friends": {
			friendCount: 10,
			expErr:      nil,
		},
	} {
		t.Run(s, func(t *testing.T) {

			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				userOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.email,
				}
				userOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
				var expOut model.Users
				for i := 1; i <= tc.friendCount; i++ {
					friendEmail := fmt.Sprintf("friend%d@test.com", i)
					friendOrm := &orm.User{
						UserID:    0,
						UserEmail: friendEmail,
					}
					friendOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

					relaUtoFOrm := &orm.Relationship{
						ReceiverID: userOrm.UserID,
						SenderID:   friendOrm.UserID,
						Status:     orm.StatusFriends,
					}
					relaFtoUOrm := &orm.Relationship{
						ReceiverID: friendOrm.UserID,
						SenderID:   userOrm.UserID,
						Status:     orm.StatusFriends,
					}

					relaUtoFOrm.Insert(context.Background(), executor, boil.Blacklist())
					relaFtoUOrm.Insert(context.Background(), executor, boil.Blacklist())

					expOut = append(expOut, model.User{
						UserID:    friendOrm.UserID,
						UserEmail: friendEmail,
					})
				}

				repo := New(executor)

				// When:
				out, err := repo.GetFriends(context.Background(), model.User{
					UserID:    userOrm.UserID,
					UserEmail: testConst.email,
				})

				// Then:
				require.ElementsMatch(t, expOut, out)
				require.ErrorIs(t, err, tc.expErr)
			})
		})
	}
}
