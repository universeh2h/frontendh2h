package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/model"
	"github.com/universeh2h/report/internal/services"
	"github.com/universeh2h/report/pkg/response"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(services *services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: services,
	}
}

func (h *ProductHandler) GetAnalytics(r *fiber.Ctx) error {
	startDate := r.Query("startDate")
	endDate := r.Query("endDate")

	data, err := h.service.Analytics(r.Context(), model.PaginationParams{
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		return response.ErrorResponse(r, http.StatusInternalServerError, "Internal Server Error", err.Error())

	}

	return response.SuccessResponse(r, http.StatusOK, "Create Category Successfully", data)

}
func (h *ProductHandler) GetProductTrxTerbanyak(r *fiber.Ctx) error {
	startDate := r.Query("startDate")
	endDate := r.Query("endDate")
	kodeReseller := r.Query("kodeReseller")

	data, err := h.service.GetProductTrxTerbanyak(r.Context(), startDate, endDate, kodeReseller)

	if err != nil {
		return response.ErrorResponse(r, http.StatusInternalServerError, "Internal Server Error", err.Error())

	}

	return response.SuccessResponse(r, http.StatusOK, "Au Ahh Successfully", data)

}
func (h *ProductHandler) GetTrxTercuan(r *fiber.Ctx) error {
	startDate := r.Query("startDate")
	endDate := r.Query("endDate")
	kode_reseller := r.Query("kodeReseller")

	data, err := h.service.GetTrxTercuan(r.Context(), startDate, endDate, kode_reseller)

	if err != nil {
		return response.ErrorResponse(r, http.StatusInternalServerError, "Internal Server Error", err.Error())

	}

	return response.SuccessResponse(r, http.StatusOK, "Create Category Successfully", data)

}
