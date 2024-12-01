package testhelper

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func SetupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.Exec("CREATE TABLE loans (id UUID, borrower_id INTEGER, rate INTEGER, principal_amount INTEGER, state TEXT)")
	db.Exec("CREATE TABLE loan_approvals (loan_id UUID, employee_id INTEGER, proof TEXT)")
	return db
}
