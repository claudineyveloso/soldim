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

// func main() {
// 	cfg := configs.Config{
// 		Host:       configs.Envs.Host,
// 		Port:       configs.Envs.Port,
// 		DBUser:     configs.Envs.DBUser,
// 		DBPassword: configs.Envs.DBPassword,
// 		DBName:     configs.Envs.DBName,
// 	}
//
// 	db, err := db.NewPostgresSQLStorage(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	initStorage(db)
//
// 	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
// 	if err := server.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func initStorage(db *sql.DB) {
// 	err := db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("DB: Successfully connected!")
// }

// package main
//
// import (
//     "fmt"
//     "log"
//     "net/http"
//     "os"
// )
//
// func handler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Hello, world!")
// }
//
// func main() {
//     http.HandleFunc("/", handler)
//
//     // Use a porta fornecida pela variável de ambiente PORT
//     port := os.Getenv("PORT")
//     if port == "" {
//       port = "8080"
//       log.Fatal("Porta não definida na variável de ambiente PORT")
//     }
//
//     log.Printf("Servidor rodando em http://localhost:%s/", port)
//     log.Fatal(http.ListenAndServe(":"+port, nil))
// }
