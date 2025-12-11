package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/universeh2h/report/internal/routes"
	"github.com/universeh2h/report/pkg/config"
	loggerPkg "github.com/universeh2h/report/pkg/logger"
)

func main() {
	// Initialize logger
	log := loggerPkg.NewLogger()

	// Database configuration
	dbConfig := &config.DBCONF{
		Host:     getEnv("DB_HOST", "172.29.64.41"),
		Username: getEnv("DB_USERNAME", "client12"),
		Password: getEnv("DB_PASSWORD", "Kecapasin123+"),
		DB:       getEnv("DB_NAME", "otomax"),
		Port:     getEnv("DB_PORT", "1433"),
	}

	// Connect to database
	db, err := dbConfig.NewDatabaseConnection()
	if err != nil {
		log.Logger.Fatal("Failed to connect to database: ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Logger.Error("Failed to close database connection: ", err)
		}
	}()

	log.Logger.Info("Database connected successfully")

	app := fiber.New(fiber.Config{
		AppName:      "Report API v1.0",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Logger.Error(fmt.Sprintf("Error %d: %s", code, err.Error()))

			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
				"code":    code,
			})
		},
	})

	// Middleware setup
	setupMiddleware(app)

	// Setup routes
	routes.SetupRoutes(app, db)

	// Graceful shutdown
	go func() {
		port := getEnv("PORT", "4000")
		log.Logger.Info(fmt.Sprintf("Server starting on port %s", port))

		if err := app.Listen(":" + port); err != nil {
			log.Logger.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Logger.Error("Server forced to shutdown: ", err)
	}

	log.Logger.Info("Server exited")
}

func setupMiddleware(app *fiber.App) {
	// Recovery middleware
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} - ${latency}\n",
		Output: os.Stdout,
	}))

	// CORS middleware
	app.Use(cors.New(setupCORS()))
}

func setupCORS() cors.Config {
	env := getEnv("APP_ENV", "development")

	config := cors.Config{
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		ExposeHeaders:    "Content-Length,Content-Type",
		AllowCredentials: true,
	}

	if env == "production" {
		// Production: Strict CORS
		config.AllowOrigins = "https://yourdomain.com,https://www.yourdomain.com"
	} else {
		// Development: More permissive CORS
		allowedOrigins := getEnv("ALLOWED_ORIGINS", "")
		if allowedOrigins != "" {
			// Use environment variable if set
			config.AllowOrigins = allowedOrigins
		} else {
			config.AllowOrigins = "http://localhost:5173,http://localhost:3000,https://pf69lscd-5173.asse.devtunnels.ms,https://mvbdbbk3-5173.asse.devtunnels.ms,http://103.184.122.173:5173"
		}
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
