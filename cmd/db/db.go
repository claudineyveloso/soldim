package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	_ "github.com/lib/pq"
)

func NewPostgresSQLStorage(cfg configs.Config) (*sql.DB, error) {
	sslMode := "disable"
	if cfg.Environment == "prod" {
		sslMode = "require" // Use SSL in production
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return dbConn, err
}
