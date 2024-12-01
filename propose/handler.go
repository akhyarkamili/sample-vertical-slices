package propose

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	cmd *Command
}

func NewHandler(app *Command) *Handler {
	return &Handler{cmd: app}
}

type httpRequest struct {
	BorrowerID      int `json:"borrower_id"`
	Rate            int `json:"rate"`
	PrincipalAmount int `json:"principal_amount"`
}

func (h *Handler) Handle(ctx echo.Context) error {
	var req httpRequest
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	request := request{
		BorrowerID:      req.BorrowerID,
		Rate:            req.Rate,
		PrincipalAmount: req.PrincipalAmount,
	}

	if err := h.cmd.Propose(request); err != nil {
		return ctx.JSON(500, map[string]string{"error": "internal server error"})
	}

	return ctx.JSON(200, map[string]string{"status": "ok"})
}
