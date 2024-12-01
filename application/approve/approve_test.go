package approve

import (
	"loan-management/application/common"
	"loan-management/application/propose"
	"loan-management/testhelper"
	"testing"

	"github.com/google/uuid"
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
			Proof:      "https://google.com",
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
		assert.Equal(t, "https://google.com", proof)
	})

	t.Run("command refuses invalid request", func(t *testing.T) {
		// Arrange
		approveCmd := NewCommand(nil)

		tests := []struct {
			name    string
			request Request
		}{
			{
				"empty request",
				Request{},
			},
			{
				"invalid url",
				Request{
					LoanID:     uuid.New(),
					EmployeeID: 1,
					Proof:      "abcde",
				},
			},
			{
				"invalid ID",
				Request{
					LoanID:     uuid.Nil,
					EmployeeID: 1,
					Proof:      "https://google.com",
				},
			},
			{
				"invalid employee ID",
				Request{
					LoanID:     uuid.New(),
					EmployeeID: 0,
					Proof:      "https://google.com",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				err := approveCmd.Approve(tt.request)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, common.ErrInvalidRequest)
			})
		}
	})
}
