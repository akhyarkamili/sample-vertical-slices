package testhelper

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.Exec("CREATE TABLE loans (id UUID, borrower_id INTEGER, rate INTEGER, principal_amount INTEGER, state TEXT)")
	return db
}
