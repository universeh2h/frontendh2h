package main

import (
	"fmt"
	"net/http"

	"github.com/universeh2h/report/internal/handler"
	"github.com/universeh2h/report/internal/repositories"
	"github.com/universeh2h/report/internal/services"
	"github.com/universeh2h/report/pkg/config"
)

func main() {

	// Database configuration
	dbConfig := &config.DBCONF{
		Host:     "172.29.64.41",
		Username: "client12",
		Password: "Kecapasin123+",
		DB:       "otomax",
		Port:     "1433",
	}

	// Connect to database
	db, err := dbConfig.NewDatabaseConnection()
	if err != nil {
		fmt.Printf("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("failed to connect database")
		}
	}()

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionsService(transactionRepo)
	transactionService.Start()

	handler := handler.NewTransactionHandler(transactionService)

	fmt.Printf("application is running")
	http.HandleFunc("/ws/transactions", handler.CheckTransactionsRealTime)
	http.ListenAndServe(":1000", nil)

}
