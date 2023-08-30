package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
	"github.com/hxcuber/friends-management/api/pkg/testutil"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"testing"
)

const (
	host     = "localhost"
	port     = 5432
	username = "hxcuber"
	password = "hxcuber"
	dbname   = "friends"
)

func TestImpl_CreateRelationship(t *testing.T) {
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
			duplicate: false,
			expErr:    nil,
		},
		"duplicate": {
			duplicate: true,
			expErr: errors.WithMessage(&pq.Error{
				Severity:         "ERROR",
				Code:             "23505",
				Message:          "duplicate key value violates unique constraint \"relationships_pkey\"",
				Detail:           "Key (receiver_id, sender_id):(63, 64) already exists.",
				Hint:             "",
				Position:         "",
				InternalPosition: "",
				InternalQuery:    "",
				Where:            "",
				Schema:           "public",
				Table:            "relationships",
				Column:           "",
				DataTypeName:     "",
				Constraint:       "relationships_pkey",
				File:             "nbtinsert.c",
				Line:             "671",
				Routine:          "_bt_check_unique"},
				"orm: unable to insert into relationships"),
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
				relaModel := model.Relationship{
					ReceiverID: receiverOrm.UserID,
					SenderID:   senderOrm.UserID,
					Status:     testConst.status,
				}
				repo := New(executor)

				// When:
				err := repo.CreateRelationship(context.Background(), relaModel)

				// Then:
				if tc.expErr != nil {
					require.ErrorContains(t, err, tc.expErr.Error())
				} else {
					require.ErrorIs(t, err, tc.expErr)
				}
			})
		})
	}
}
