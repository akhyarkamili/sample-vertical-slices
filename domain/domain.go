package domain

import "github.com/shopspring/decimal"

type state int

const (
	stateProposed state = iota
	stateApproved
	stateInvested
)

func (ls state) String() string {
	switch ls {
	case stateProposed:
		return "proposed"
	case stateApproved:
		return "approved"
	case stateInvested:
		return "invested"
	}
	return "unknown"
}

func fromString(s string) state {
	switch s {
	case "proposed":
		return stateProposed
	case "approved":
		return stateApproved
	case "invested":
		return stateInvested
	}
	return -1
}

type Approval struct {
	Proof      string
	EmployeeID int
}

type Investment struct {
	amount     decimal.Decimal
	investorID int
}

type Loan struct {
	BorrowerID      int
	Rate            int
	PrincipalAmount decimal.Decimal

	Investments []Investment
	Approval    *Approval

	state state
}

func NewLoan(borrowerID, rate int, principalAmount int) *Loan {
	return &Loan{
		BorrowerID:      borrowerID,
		Rate:            rate,
		PrincipalAmount: decimal.Decimal(decimal.NewFromInt(int64(principalAmount))),
		state:           stateProposed,
	}
}

func Load(borrowerID, rate, principalAmount int, state string, approval *Approval, investments []Investment) Loan {
	return Loan{
		BorrowerID:      borrowerID,
		Rate:            rate,
		PrincipalAmount: decimal.NewFromInt(int64(principalAmount)),
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

func (l *Loan) Invest(amount decimal.Decimal, investorID int) (overflow decimal.Decimal) {
	if l.state != stateApproved {
		return
	}

	available := l.PrincipalAmount.Sub(l.TotalInvested())
	investmentAmount := decimal.Min(available, amount)

	i := Investment{amount: investmentAmount, investorID: investorID}
	l.Investments = append(l.Investments, i)

	if l.TotalInvested().Equal(l.PrincipalAmount) {
		l.MarkInvested()
	}

	return amount.Sub(investmentAmount)
}

func (l *Loan) TotalInvested() decimal.Decimal {
	var total decimal.Decimal
	for _, inv := range l.Investments {
		total = total.Add(inv.amount)
	}
	return total
}

func (l *Loan) MarkInvested() {
	if l.state != stateApproved {
		return
	}
	l.state = stateInvested
}
