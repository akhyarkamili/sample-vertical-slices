package invest

import (
	"loan-management/application/common"
	"loan-management/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Repository interface {
	Get(loanID uuid.UUID) (domain.Loan, error)
	SaveInvestments(id uuid.UUID, loan domain.Loan) error
}

type Command struct {
	repo Repository
}

func NewCommand(repo Repository) Command {
	return Command{
		repo: repo,
	}
}

type Request struct {
	LoanID     uuid.UUID       `json:"id" validate:"required"`
	InvestorID int             `json:"investor_id" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required"`
}

func (r *Request) Validate() error {
	return validator.New().Struct(r)
}

func (cmd *Command) Invest(req Request) error {
	if err := req.Validate(); err != nil {
		return common.ErrInvalidRequest
	}

	loan, err := cmd.repo.Get(req.LoanID)
	if err != nil {
		return err
	}

	loan.Invest(req.Amount, req.InvestorID)
	return cmd.repo.SaveInvestments(req.LoanID, loan)
}
