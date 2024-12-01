package domain

import "github.com/shopspring/decimal"

type Rupiah decimal.Decimal

type state int

const (
	stateProposed state = iota
	stateApproved
)

func (ls state) String() string {
	switch ls {
	case stateProposed:
		return "proposed"
	case stateApproved:
		return "approved"
	}
	return "unknown"
}

type Loan struct {
	BorrowerID      int
	Rate            int
	PrincipalAmount Rupiah

	state state
}

func NewLoan(borrowerID, rate int, principalAmount int) *Loan {
	return &Loan{
		BorrowerID:      borrowerID,
		Rate:            rate,
		PrincipalAmount: Rupiah(decimal.NewFromInt(int64(principalAmount))),
		state:           stateProposed,
	}
}

func (l *Loan) State() string {
	return l.state.String()
}

func (l *Loan) approve() {
	if l.state != stateProposed {
		return
	}
	l.state = stateApproved
}
