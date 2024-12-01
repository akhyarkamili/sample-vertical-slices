package propose

import (
	"loan-management/domain"

	"github.com/google/uuid"
)

type Repository interface {
	Create(id uuid.UUID, loan domain.Loan) error
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
	BorrowerID      int
	Rate            int
	PrincipalAmount int
}

func (r *Request) Validate() error {
	if r.BorrowerID <= 0 {
		return ErrInvalidBorrowerID
	}
	if r.Rate <= 0 {
		return ErrInvalidRate
	}
	if r.PrincipalAmount <= 0 {
		return ErrInvalidPrincipalAmount
	}
	return nil
}

func (ps *Command) Propose(request Request) (loanID uuid.UUID, err error) {
	if err := request.Validate(); err != nil {
		return uuid.Nil, err
	}
	loan := domain.NewLoan(request.BorrowerID, request.Rate, request.PrincipalAmount)

	loanID = uuid.New()
	err = ps.repo.Create(loanID, *loan)
	if err != nil {
		return uuid.Nil, err
	}

	return loanID, nil
}
