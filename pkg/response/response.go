package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(response)
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err string) error {
	response := Response{
		Success: false,
		Message: message,
		Error:   err,
	}
	return c.Status(statusCode).JSON(response)
}
