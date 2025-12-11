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
	handlers := handler.NewProductHandler(service)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionServices := services.NewTransactionsService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionServices)

	moduleRepo := repositories.NewModulOtomax(db)
	moduleService := services.NewModulService(moduleRepo)
	moduleHandler := handler.NewModulHandler(moduleService)

	api := r.Group("/api/v1")
	api.Get("", handlers.GetAnalytics)
	api.Get("/trxtercuan", handlers.GetTrxTercuan)
	api.Get("/trxterbanyak", handlers.GetProductTrxTerbanyak)

	api.Get("/transactions", transactionHandler.CheckTransactionsRealTime)
	api.Get("/report", handlers.Report)
	api.Get("/modul-otomax", moduleHandler.GetAllModulOtomax)
}
