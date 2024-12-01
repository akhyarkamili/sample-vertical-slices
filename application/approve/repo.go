package approve

import (
	"database/sql"
	"loan-management/domain"

	"github.com/google/uuid"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r repo) Get(id uuid.UUID) (domain.Loan, error) {
	rows, err := r.db.Query("SELECT borrower_id, rate, principal_amount, state FROM loans WHERE id = ?", id)
	if err != nil {
		return domain.Loan{}, err
	}
	defer rows.Close()
	var borrowerID, rate, principalAmount int
	var state string
	for rows.Next() {
		err := rows.Scan(&borrowerID, &rate, &principalAmount, &state)
		if err != nil {
			return domain.Loan{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return domain.Loan{}, err
	}

	return domain.Load(borrowerID, rate, principalAmount, state, nil), nil
}

func (r repo) SaveApproval(id uuid.UUID, loan domain.Loan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO loan_approvals (loan_id, employee_id, proof) VALUES (?, ?, ?)", id, loan.Approval.EmployeeID, loan.Approval.Proof)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE loans SET state = ? WHERE id = ?", loan.State(), id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}
