package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	_ "github.com/lib/pq"
)

func NewPostgresSQLStorageAAA(cfg configs.Config) (*sql.DB, error) {
	var connStr string

	// Verifica se a DATABASE_URL está definida no ambiente
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// Se DATABASE_URL não estiver definida, constrói a string manualmente
		sslMode := "disable"
		if cfg.Environment == "prod" {
			sslMode = "require" // Use SSL in production
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verifica se a conexão está funcionando
	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbConn, nil
}

func NewPostgresSQLStorage(cfg configs.Config) (*sql.DB, error) {
	sslMode := "disable"
	if cfg.Environment == "prod" {
		sslMode = "require" // Use SSL in production
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)

	// connStr := "host=c3gtj1dt5vh48j.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com port=5432 user=u8v8gfju2nbfqn password=p009827e7e08f28b44b9ba56751f82543345f445a718fb44624dd5b653e0238bd dbname=ddcvr5rele9132 sslmode=require"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return dbConn, err
}
