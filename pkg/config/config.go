package config

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

type DBCONF struct {
	Host     string
	Username string
	DB       string
	Password string
	Port     string
}

// NewDatabaseConnection creates and returns a new database connection
func (conf *DBCONF) NewDatabaseConnection() (*sql.DB, error) {
	// Format connection string for SQL Server
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		conf.Host, conf.Username, conf.Password, conf.Port, conf.DB)

	// Open database connection
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}

// Alternative connection string format with additional parameters
func (conf *DBCONF) NewDatabaseConnectionWithParams() (*sql.DB, error) {
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.DB)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
