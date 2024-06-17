package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	_ "github.com/lib/pq"
)

func NewPostgresSQLStorage(cfg configs.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PublicHost, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return dbConn, err
}
