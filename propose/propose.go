package propose

import (
	"loan-management/domain"
)

type Repository interface {
	Save(loan domain.Loan) error
}

type Command struct {
	repo Repository
}

func NewCommand(repo Repository) *Command {
	return &Command{
		repo: repo,
	}
}

type request struct {
	BorrowerID      int
	Rate            int
	PrincipalAmount int
}

func (r *request) Validate() error {
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

func (ps *Command) Propose(request request) error {
	if err := request.Validate(); err != nil {
		return err
	}
	loan := domain.NewLoan(request.BorrowerID, request.Rate, request.PrincipalAmount)
	return ps.repo.Save(*loan)
}
