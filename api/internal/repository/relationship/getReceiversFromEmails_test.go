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

func TestImpl_GetReceiversFromEmails(t *testing.T) {
	testConst := struct {
		senderEmail string
	}{
		senderEmail: "sender@test.com",
	}
	type test struct {
		subscribedCount int
		blockedCount    int
		friendCount     int
		noneCount       int
		absentCount     int
		expErr          error
	}
	os.Setenv("DB_URL", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname))
	for s, tc := range map[string]test{
		"empty": {
			subscribedCount: 0,
			blockedCount:    0,
			friendCount:     0,
			noneCount:       0,
			absentCount:     0,
			expErr:          nil,
		},
		"test_1": {
			subscribedCount: 3,
			blockedCount:    7,
			friendCount:     8,
			noneCount:       7,
			absentCount:     9,
			expErr:          nil,
		},
		"test_2": {
			subscribedCount: 8,
			blockedCount:    4,
			friendCount:     5,
			noneCount:       10,
			absentCount:     3,
			expErr:          nil,
		},
		"test_3": {
			subscribedCount: 1,
			blockedCount:    6,
			friendCount:     5,
			noneCount:       9,
			absentCount:     8,
			expErr:          nil,
		},
		"test_4": {
			subscribedCount: 0,
			blockedCount:    5,
			friendCount:     10,
			noneCount:       4,
			absentCount:     9,
			expErr:          nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				senderOrm := &orm.User{
					UserID:    0,
					UserEmail: testConst.senderEmail,
				}
				senderOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

				var emails []string
				var expOut model.Users
				for i := 1; i <= tc.subscribedCount; i++ {
					subscriberEmail := fmt.Sprintf("subscriber%d@test.com", i)

					subscriberOrm := &orm.User{
						UserID:    0,
						UserEmail: subscriberEmail,
					}
					subscriberOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

					relaOrm := &orm.Relationship{
						ReceiverID: subscriberOrm.UserID,
						SenderID:   senderOrm.UserID,
						Status:     orm.StatusRSubscribedS,
					}
					relaOrm.Insert(context.Background(), executor, boil.Blacklist())

					emails = append(emails, subscriberEmail)
					expOut = append(expOut, model.User{
						UserID:    subscriberOrm.UserID,
						UserEmail: subscriberEmail,
					})
				}
				for i := 1; i <= tc.blockedCount; i++ {
					blockedEmail := fmt.Sprintf("blocker%d@test.com", i)

					blockedOrm := &orm.User{
						UserID:    0,
						UserEmail: blockedEmail,
					}
					blockedOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

					relaOrm := &orm.Relationship{
						ReceiverID: blockedOrm.UserID,
						SenderID:   senderOrm.UserID,
						Status:     orm.StatusRBlockedS,
					}
					relaOrm.Insert(context.Background(), executor, boil.Blacklist())

					emails = append(emails, blockedEmail)
				}
				for i := 1; i <= tc.friendCount; i++ {
					friendEmail := fmt.Sprintf("friend%d@test.com", i)
					friendOrm := &orm.User{
						UserID:    0,
						UserEmail: friendEmail,
					}
					friendOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

					relaUtoFOrm := &orm.Relationship{
						ReceiverID: senderOrm.UserID,
						SenderID:   friendOrm.UserID,
						Status:     orm.StatusFriends,
					}
					relaFtoUOrm := &orm.Relationship{
						ReceiverID: friendOrm.UserID,
						SenderID:   senderOrm.UserID,
						Status:     orm.StatusFriends,
					}

					relaUtoFOrm.Insert(context.Background(), executor, boil.Blacklist())
					relaFtoUOrm.Insert(context.Background(), executor, boil.Blacklist())

					emails = append(emails, friendEmail)
					expOut = append(expOut, model.User{
						UserID:    friendOrm.UserID,
						UserEmail: friendEmail,
					})
				}
				for i := 1; i <= tc.noneCount; i++ {
					noneEmail := fmt.Sprintf("none%d@test.com", i)

					noneOrm := &orm.User{
						UserID:    0,
						UserEmail: noneEmail,
					}
					noneOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))

					emails = append(emails, noneEmail)
					expOut = append(expOut, model.User{
						UserID:    noneOrm.UserID,
						UserEmail: noneEmail,
					})
				}
				for i := 1; i <= tc.absentCount; i++ {
					absentEmail := fmt.Sprintf("absent%d@test.com", i)

					emails = append(emails, absentEmail)
				}
				repo := New(executor)

				// When:
				out, err := repo.GetReceiversFromEmails(context.Background(), model.User{
					UserID:    senderOrm.UserID,
					UserEmail: testConst.senderEmail,
				}, emails)

				// Then:
				require.ElementsMatch(t, expOut, out)
				require.ErrorIs(t, err, tc.expErr)
			})
		})
	}
}
