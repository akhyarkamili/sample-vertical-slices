package loan_query

import (
	"database/sql"
	"errors"
	"fmt"
	"loan-management/application/common"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type Query struct {
	ID uuid.UUID `param:"loan_id" validate:"required,uuid"`
}

type QueryResult struct {
	LoanID          uuid.UUID       `json:"loan_id"`
	PrincipalAmount decimal.Decimal `json:"principal_amount"`
	State           string          `json:"state"`
}

func (h *Handler) Handle(ctx echo.Context) error {
	var req Query
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	result, err := h.query(req)
	if err != nil {
		if errors.Is(err, common.ErrInvalidRequest) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		fmt.Printf("Error: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return ctx.JSON(http.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) query(req Query) (result QueryResult, err error) {
	rows, err := h.db.Query("SELECT id as loan_id, principal_amount, state FROM loans WHERE id = ?", req.ID)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&result.LoanID, &result.PrincipalAmount, &result.State)
		if err != nil {
			return result, err
		}
	}
	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil

}
