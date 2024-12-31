package invest

import (
	"database/sql"
	"loan-management/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Get(loanID uuid.UUID) (domain.Loan, error) {
	var borrowerID, rate int
	var principalAmount int
	var state string
	var investments []domain.Investment

	// Fetch loan details
	err := r.db.QueryRow("SELECT borrower_id, rate, principal_amount, state FROM loans WHERE id = ?", loanID).
		Scan(&borrowerID, &rate, &principalAmount, &state)
	if err != nil {
		return domain.Loan{}, err
	}

	// Fetch investments
	rows, err := r.db.Query("SELECT amount, investor_id FROM investments WHERE loan_id = ?", loanID)
	if err != nil {
		return domain.Loan{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var amount decimal.Decimal
		var investorID int
		err := rows.Scan(&amount, &investorID)
		if err != nil {
			return domain.Loan{}, err
		}
		investments = append(investments, domain.Investment{Amount: amount, InvestorID: investorID})
	}

	if err := rows.Err(); err != nil {
		return domain.Loan{}, err
	}

	return domain.Load(borrowerID, rate, principalAmount, state, nil, investments), nil
}

func (r *repo) SaveInvestments(loanID uuid.UUID, loan domain.Loan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insert investments
	for _, investment := range loan.Investments {
		_, err := tx.Exec("INSERT INTO investments (loan_id, amount, investor_id) VALUES (?, ?, ?)",
			loanID, investment.Amount, investment.InvestorID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update loan state
	_, err = tx.Exec("UPDATE loans SET state = ? WHERE id = ?", loan.State(), loanID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
