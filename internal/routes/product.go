package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/internal/handler"
	"github.com/universeh2h/report/internal/repositories"
	"github.com/universeh2h/report/internal/services"
)

func SetupRoutes(r *fiber.App, db *sql.DB) {
	repo := repositories.NewProductRepository(db)
	service := services.NewProductServices(repo)
	handler := handler.NewProductHandler(service)

	api := r.Group("/api/v1")
	api.Get("", handler.GetAnalytics)
	api.Get("/trxtercuan", handler.GetTrxTercuan)
	api.Get("/trxterbanyak", handler.GetProductTrxTerbanyak)
}
