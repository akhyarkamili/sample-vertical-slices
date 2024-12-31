package invest

import (
	"errors"
	"loan-management/application/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	cmd Command
}

func NewHandler(cmd Command) *Handler {
	return &Handler{
		cmd: cmd,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.cmd.Invest(req)
	if err != nil {
		if errors.Is(err, common.ErrInvalidRequest) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return ctx.JSON(http.StatusOK, map[string]any{"success": true})
}
