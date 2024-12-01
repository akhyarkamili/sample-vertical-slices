package approve

import (
	"loan-management/propose"
	"loan-management/testhelper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprove(t *testing.T) {
	t.Run("command approves a proposed loan", func(t *testing.T) {
		// Arrange
		db := testhelper.SetupDB(t)
		proposeRepo := propose.NewRepository(db)
		proposeCmd := propose.NewCommand(proposeRepo)
		id, err := proposeCmd.Propose(propose.Request{
			BorrowerID:      1,
			Rate:            10,
			PrincipalAmount: 1000000,
		})
		require.NoError(t, err)

		repo := NewRepository(db)
		approveCmd := NewCommand(repo)
		request := Request{
			LoanID:     id,
			EmployeeID: 1,
			Proof:      "http://google.com",
		}
		// Act
		err = approveCmd.Approve(request)

		// Assert
		assert.NoError(t, err)
		rows, err := db.Query("SELECT state FROM loans WHERE id = ?", id)
		require.NoError(t, err)
		defer rows.Close()
		var state string
		for rows.Next() {
			err := rows.Scan(&state)
			require.NoError(t, err)
		}
		require.NoError(t, rows.Err())
		assert.Equal(t, "approved", state)

		rows, err = db.Query("SELECT employee_id, proof FROM loan_approvals WHERE loan_id = ?", id)
		require.NoError(t, err)
		defer rows.Close()
		var employeeID int
		var proof string
		for rows.Next() {
			err := rows.Scan(&employeeID, &proof)
			require.NoError(t, err)
		}
		require.NoError(t, rows.Err())
		assert.Equal(t, 1, employeeID)
		assert.Equal(t, "http://google.com", proof)
	})
}
