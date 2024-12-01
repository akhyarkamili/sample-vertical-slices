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

func fromString(s string) state {
	switch s {
	case "proposed":
		return stateProposed
	case "approved":
		return stateApproved
	}
	return -1
}

type Approval struct {
	Proof      string
	EmployeeID int
}

type Loan struct {
	BorrowerID      int
	Rate            int
	PrincipalAmount Rupiah

	Approval *Approval

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

func Load(borrowerID, rate, principalAmount int, state string, approval *Approval) Loan {
	return Loan{
		BorrowerID:      borrowerID,
		Rate:            rate,
		PrincipalAmount: Rupiah(decimal.NewFromInt(int64(principalAmount))),
		state:           fromString(state),
		Approval:        approval,
	}
}

func (l *Loan) State() string {
	return l.state.String()
}

func (l *Loan) Approve(proof string, employeeID int) {
	if l.state != stateProposed {
		return
	}
	l.state = stateApproved
	l.Approval = &Approval{
		Proof:      proof,
		EmployeeID: employeeID,
	}
}
