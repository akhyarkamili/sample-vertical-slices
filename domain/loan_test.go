package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLoan_Invest(t *testing.T) {
	t.Run("valid investment", func(t *testing.T) {
		loan := NewLoan(1, 10, 1000000)
		loan.Approve("https://google.com/valid_proof", 1)

		loan.Invest(decimal.NewFromInt32(1000), 0)

		assert.Equal(t, loan.TotalInvested(), decimal.NewFromInt32(1000))
	})

	t.Run("overflowing investment", func(t *testing.T) {
		loan := NewLoan(1, 10, 1000000)
		loan.Approve("https://google.com/valid_proof", 1)

		leftover := loan.Invest(decimal.NewFromInt32(1000001), 0)

		assert.Equal(t, loan.TotalInvested(), loan.PrincipalAmount)
		assert.Equal(t, decimal.NewFromInt32(1), leftover)
	})

	t.Run("multiple investment", func(t *testing.T) {
		loan := NewLoan(1, 10, 1000000)
		loan.Approve("https://google.com/valid_proof", 1)

		_ = loan.Invest(decimal.NewFromInt32(999999), 0)
		leftover := loan.Invest(decimal.NewFromInt32(3), 0)

		assert.Equal(t, loan.TotalInvested(), loan.PrincipalAmount)
		assert.Equal(t, decimal.NewFromInt32(2), leftover)
	})
}
