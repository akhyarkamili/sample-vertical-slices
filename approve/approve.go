package approve

import (
	"loan-management/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Repository interface {
	Get(loanID uuid.UUID) (domain.Loan, error)
	SaveApproval(id uuid.UUID, loan domain.Loan) error
}

type Command struct {
	repo Repository
}

func NewCommand(repo Repository) *Command {
	return &Command{
		repo: repo,
	}
}

type Request struct {
	LoanID     uuid.UUID `validate:"required"`
	EmployeeID int       `validate:"required"`
	Proof      string    `validate:"required,url"`
}

func (r *Request) Validate() error {
	return validator.New().Struct(r)
}

func (ps *Command) Approve(req Request) error {
	if err := req.Validate(); err != nil {
		return ErrInvalidRequest
	}

	loan, err := ps.repo.Get(req.LoanID)
	if err != nil {
		return err
	}

	loan.Approve(req.Proof, req.EmployeeID)
	return ps.repo.SaveApproval(req.LoanID, loan)
}
