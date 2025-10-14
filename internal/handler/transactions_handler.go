package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/services"
	"github.com/universeh2h/report/pkg/response"
)

type TransactionHandler struct {
	service *services.TransactionsService
}
type CustomError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewTransactionHandler(service *services.TransactionsService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

// CheckTransactionsRealTime handle WebSocket connection
func (h *TransactionHandler) CheckTransactionsRealTime(c *fiber.Ctx) error {
	// Set CORS headers if needed

	data, err := h.service.GetTransactions(context.Background())
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "failed to get transactions", err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "transactions successfully", data)
}
