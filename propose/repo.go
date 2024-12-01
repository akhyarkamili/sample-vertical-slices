package propose

import (
	"database/sql"
	"loan-management/domain"

	"github.com/shopspring/decimal"
)

type loanModel struct {
	BorrowerID      int
	Rate            int
	PrincipalAmount decimal.Decimal
	State           string
}

func FromDomain(loan domain.Loan) loanModel {
	return loanModel{
		BorrowerID:      loan.BorrowerID,
		Rate:            loan.Rate,
		PrincipalAmount: decimal.Decimal(loan.PrincipalAmount),
		State:           loan.State(),
	}
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Save(l domain.Loan) error {
	loan := FromDomain(l)
	_, err := r.db.Exec("INSERT INTO loans (borrower_id, rate, principal_amount, state) VALUES (?, ?, ?, ?)", loan.BorrowerID, loan.Rate, loan.PrincipalAmount, loan.State)
	return err
}
