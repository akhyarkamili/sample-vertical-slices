package propose

import (
	"loan-management/domain"
	"loan-management/testhelper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

type mockRepository struct{}

func (mr *mockRepository) Save(_ uuid.UUID, _ domain.Loan) error {
	return nil
}

func TestPropose(t *testing.T) {
	t.Run("propose exists", func(t *testing.T) {
		repo := &mockRepository{}
		service := NewCommand(repo)
		_, err := service.Propose(request{
			BorrowerID:      1,
			Rate:            10,
			PrincipalAmount: 1000000,
		})
		assert.NoError(t, err)
	})

	t.Run("command saves the proposed state and returns correct ID", func(t *testing.T) {
		db := testhelper.SetupDB(t)
		repo := NewRepository(db)
		service := NewCommand(repo)
		id, err := service.Propose(request{
			BorrowerID:      1,
			Rate:            10,
			PrincipalAmount: 1000000,
		})
		assert.NoError(t, err)

		// query from DB
		rows, err := db.Query("SELECT borrower_id, rate, principal_amount, state FROM loans WHERE id = ?", id)
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
