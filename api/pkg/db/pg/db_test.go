package pg

import (
	"context"
	"database/sql"
	"github.com/hxcuber/friends-management/api/pkg/env"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFriendsDB(t *testing.T) {
	pool, err := sql.Open("postgres", env.GetAndValidateF("DB_URL"))
	require.NoError(t, err)

	db := &friendsDB{DB: pool}

	_, err = db.Exec("DROP TABLE IF EXISTS instrumented_test_transactions")
	require.NoError(t, err)
	_, err = db.Exec("CREATE TABLE instrumented_test_transactions (id SERIAL PRIMARY KEY)")
	require.NoError(t, err)
	defer func() {
		_, err = db.Exec("DROP TABLE IF EXISTS instrumented_test_transactions")
		require.NoError(t, err)
	}()

	_, err = db.QueryContext(context.Background(), "SELECT * FROM instrumented_test_transactions")
	require.NoError(t, err)

	var p interface{}
	row := db.QueryRowContext(context.Background(), "SELECT * FROM instrumented_test_transactions")
	err = row.Scan(&p)
	require.Equal(t, sql.ErrNoRows, err)

	row = db.QueryRow("SELECT * FROM instrumented_test_transactions")
	err = row.Scan(&p)
	require.Equal(t, sql.ErrNoRows, err)

	_, err = db.ExecContext(context.Background(), "SELECT * FROM instrumented_test_transactions")
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)
	require.NoError(t, tx.Rollback())
}
