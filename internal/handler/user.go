package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/middleware"
	"github.com/universeh2h/report/internal/model"
	"github.com/universeh2h/report/internal/services"
	"github.com/universeh2h/report/pkg/response"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(services *services.UserService) *UserHandler {
	return &UserHandler{
		service: services,
	}
}

func (h *UserHandler) Login(a *fiber.Ctx) error {
	var input model.Login
	if err := a.BodyParser(&input); err != nil {
		return response.ErrorResponse(a, http.StatusBadRequest, "Invalid input", err.Error())
	}
	user, token, err := h.service.Login(a.Context(), input)
	if err != nil {
		fmt.Printf("errr : %s", err.Error())
		return response.ErrorResponse(a, http.StatusInternalServerError, "Failed to create category", err.Error())
	}
	middleware.NewAuthHelpers().SetAccessTokenCookie(a, token)

	return response.SuccessResponse(a, http.StatusOK, "Created successfully", user)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	username := c.Locals("username")
	if username == nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Invalid user data", "")
	}

	return response.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", usernameStr)
}
