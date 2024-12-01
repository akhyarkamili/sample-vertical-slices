package propose

import (
	"encoding/json"

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
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid Request"})
	}

	request := Request{
		BorrowerID:      req.BorrowerID,
		Rate:            req.Rate,
		PrincipalAmount: req.PrincipalAmount,
	}

	id, err := h.cmd.Propose(request)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": "internal server error"})
	}

	return ctx.JSON(200, map[string]any{
		"success": true,
		"id":      id.String(),
	})
}
