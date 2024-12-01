package approve

import (
	"encoding/json"
	"errors"
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
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.cmd.Approve(req)
	if err != nil {
		if errors.Is(err, ErrInvalidRequest) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return ctx.JSON(http.StatusOK, map[string]any{"success": true})
}
