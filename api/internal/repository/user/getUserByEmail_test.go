package user

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
	"github.com/hxcuber/friends-management/api/pkg/testutil"
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

func TestImpl_GetUserByEmail(t *testing.T) {
	input := "test@test.com"
	type test struct {
		create bool
		expOut model.User
		expErr error
	}

	os.Setenv("DB_URL", "postgres://friends_management:@pg:5432/friends_management?sslmode=disable")
	for s, tc := range map[string]test{
		"emailNotFound": {
			create: false,
			expOut: model.User{},
			expErr: ErrEmailNotFound,
		},
		"success": {
			create: true,
			expOut: model.User{UserEmail: "test@test.com"},
			expErr: nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(executor pg.BeginnerExecutor) {
				// Given:
				if tc.create {
					userOrm := &orm.User{
						UserID:    0,
						UserEmail: input,
					}
					userOrm.Insert(context.Background(), executor, boil.Blacklist(orm.UserColumns.UserID))
					tc.expOut.UserID = userOrm.UserID
				}

				repo := New(executor)

				// When:
				out, err := repo.GetUserByEmail(context.Background(), input)

				// Then:
				require.ErrorIs(t, err, tc.expErr)
				require.Equal(t, tc.expOut, out)
			})
		})
	}
}
