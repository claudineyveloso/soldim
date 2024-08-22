package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/claudineyveloso/soldim.git/cmd/api"
	"github.com/claudineyveloso/soldim.git/cmd/db"
	"github.com/claudineyveloso/soldim.git/internal/configs"
)

func main() {
	// Carrega as configurações do ambiente
	cfg := configs.Envs

	// Conecta ao banco de dados
	dbConn, err := db.NewPostgresSQLStorage(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Inicializa o banco de dados (verificação de conexão)
	initStorage(dbConn)

	// Obtém a porta da variável de ambiente
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Define a porta padrão se não for encontrada
	}

	// Inicializa o servidor da API
	server := api.NewAPIServer(fmt.Sprintf(":%s", port), dbConn)
	if err := server.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar a conexão com o banco de dados: %v", err)
	}
	log.Println("DB: Conexão estabelecida com sucesso!")
}
