package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	_ "github.com/lib/pq"
)

func NewPostgresSQLStorage(cfg configs.Config) (*sql.DB, error) {
	var connStr string

	log.Printf("Valor das variaveis de configs:%v/", cfg)

	// Verifica se DATABASE_URL está definida
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// Constrói a string de conexão manualmente
		sslMode := "disable"
		if cfg.Environment == "prod" {
			sslMode = "require" // SSL em produção
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco de dados: %w", err)
	}

	// Verifica se a conexão está funcionando
	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar o banco de dados: %w", err)
	}

	return dbConn, nil
}
