package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/handler"
	"github.com/universeh2h/report/internal/middleware"
	"github.com/universeh2h/report/internal/repositories"
	"github.com/universeh2h/report/internal/services"
)

func SetupAuthRoutes(r *fiber.App, db *sql.DB) {
	repo := repositories.NewAuthRepository(db)
	service := services.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	api := r.Group("/api/v1")
	api.Post("/login", handler.Login)

	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware())
	auth.Get("/profile", handler.GetProfile)
}
