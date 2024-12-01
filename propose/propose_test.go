package propose

import (
	"database/sql"
	"loan-management/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

type mockRepository struct{}

func (mr *mockRepository) Save(loan domain.Loan) error {
	return nil
}

func TestPropose(t *testing.T) {
	t.Run("propose exists", func(t *testing.T) {
		repo := &mockRepository{}
		service := NewCommand(repo)
		err := service.Propose(request{
			BorrowerID:      1,
			Rate:            10,
			PrincipalAmount: 1000000,
		})
		assert.NoError(t, err)
	})

	t.Run("command saves the proposed state", func(t *testing.T) {
		db := setupDB(t)
		repo := NewRepository(db)
		service := NewCommand(repo)
		err := service.Propose(request{
			BorrowerID:      1,
			Rate:            10,
			PrincipalAmount: 1000000,
		})
		assert.NoError(t, err)
		// query from DB
		rows, err := db.Query("SELECT borrower_id, rate, principal_amount, state FROM loans")
		require.NoError(t, err)
		defer rows.Close()
		var borrowerID, rate, principalAmount int
		var state string
		for rows.Next() {
			err := rows.Scan(&borrowerID, &rate, &principalAmount, &state)
			require.NoError(t, err)
		}
		require.NoError(t, rows.Err())
		assert.Equal(t, 1, borrowerID)
		assert.Equal(t, 10, rate)
		assert.Equal(t, 1000000, principalAmount)
		assert.Equal(t, "proposed", state)
	})
}

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.Exec("CREATE TABLE loans (borrower_id INTEGER, rate INTEGER, principal_amount INTEGER, state TEXT)")
	return db
}
