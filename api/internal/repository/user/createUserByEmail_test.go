package user

import (
	"context"
	"fmt"
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

func TestImpl_CreateUserByEmail(t *testing.T) {
	input := "test@test.com"
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
				Message:          "duplicate key value violates unique constraint \"users_user_email_key\"",
				Detail:           "Key (user_email)=(test@test.com) already exists.",
				Hint:             "",
				Position:         "",
				InternalPosition: "",
				InternalQuery:    "",
				Where:            "",
				Schema:           "public",
				Table:            "users",
				Column:           "",
				DataTypeName:     "",
				Constraint:       "users_user_email_key",
				Routine:          "_bt_check_unique",
			}, "orm: unable to insert into users"),
		},
	} {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				if tc.duplicate {
					userOrm := &orm.User{
						UserID:    0,
						UserEmail: input,
					}
					userOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
				}

				repo := New(executor)

				// When:
				err := repo.CreateUserByEmail(context.Background(), input)

				// Then:
				if tc.expErr != nil {
					require.ErrorContains(t, err, tc.expErr.Error())
				} else {
					require.ErrorIs(t, tc.expErr, err)
				}
			})
		})
	}
}
