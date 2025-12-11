package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/services"
	"github.com/universeh2h/report/pkg/response"
)

type ModulHandler struct {
	service *services.ModulService
}

func NewModulHandler(service *services.ModulService) *ModulHandler {
	return &ModulHandler{
		service: service,
	}
}

func (h *ModulHandler) GetAllModulOtomax(r *fiber.Ctx) error {
	date := r.Query("date")
	data, err := h.service.GetAllModulOtomax(r.Context(), date)

	if err != nil {
		return response.ErrorResponse(r, http.StatusInternalServerError, "Internal Server Error", err.Error())

	}

	return response.SuccessResponse(r, http.StatusOK, "get all modul otomax", data)
}
