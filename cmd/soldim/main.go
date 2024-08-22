package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claudineyveloso/soldim.git/cmd/api"
	"github.com/claudineyveloso/soldim.git/cmd/db"
	"github.com/claudineyveloso/soldim.git/internal/configs"
)

func main() {
	// Configuração
	cfg := configs.Envs

	// Inicialização do banco de dados
	dbConn, err := db.NewPostgresSQLStorage(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	if err := initStorage(dbConn); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

	// Inicialização do servidor
	server := api.NewAPIServer(fmt.Sprintf(":%s", cfg.Port), dbConn)
	if err := server.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func initStorage(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}
	log.Println("DB: Conectado com sucesso!")
	return nil
}
