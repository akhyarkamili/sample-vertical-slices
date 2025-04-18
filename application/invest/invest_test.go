package invest

import (
	"loan-management/application/approve"
	"loan-management/application/propose"
	"loan-management/domain"
	"loan-management/testhelper"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
}

func (m *mockRepository) Get(loanID uuid.UUID) (domain.Loan, error) {
	return domain.Load(1, 10, 1000000, "appproved", nil, []domain.Investment{}), nil
}

func (m *mockRepository) SaveInvestments(id uuid.UUID, loan domain.Loan) error {
	return nil
}

func TestInvest(t *testing.T) {
	t.Run("command invests and persists correctly", func(t *testing.T) {
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

		approveRepo := approve.NewRepository(db)
		approveCmd := approve.NewCommand(approveRepo)
		approveReq := approve.Request{
			LoanID:     id,
			EmployeeID: 1,
			Proof:      "https://google.com",
		}
		err = approveCmd.Approve(approveReq)
		require.NoError(t, err)

		// Act
		investRepo := NewRepository(db)
		investCmd := NewCommand(investRepo)
		request := Request{
			LoanID:     id,
			InvestorID: 1,
			Amount:     decimal.NewFromInt32(5000),
		}
		err = investCmd.Invest(request)
		require.NoError(t, err)

		// Assert
		var investedAmount decimal.Decimal
		err = db.QueryRow("SELECT SUM(amount) FROM investments WHERE loan_id = ?", id).Scan(&investedAmount)
		require.NoError(t, err)
		assert.Equal(t, request.Amount, investedAmount)
	})
}
