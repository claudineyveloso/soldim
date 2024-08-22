// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	"github.com/claudineyveloso/soldim.git/cmd/api"
// 	"github.com/claudineyveloso/soldim.git/cmd/db"
// 	"github.com/claudineyveloso/soldim.git/internal/configs"
// )

// func main() {
// 	cfg := configs.Config{
// 		Host:       configs.Envs.Host,
// 		Port:       configs.Envs.Port,
// 		DBUser:     configs.Envs.DBUser,
// 		DBPassword: configs.Envs.DBPassword,
// 		DBName:     configs.Envs.DBName,
// 	}

// 	db, err := db.NewPostgresSQLStorage(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	initStorage(db)

// 	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
// 	if err := server.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func initStorage(db *sql.DB) {
// 	err := db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("DB: Successfully connected!")
// }

package main

import (
	"fmt"
	"log"
	"net/http"
)

// handler para a rota "/"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Responde com "Hello, World!" para qualquer requisição GET na rota "/"
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Associa o handler helloHandler com a rota "/"
	http.HandleFunc("/", helloHandler)

	// Define a porta em que o servidor vai escutar
	port := "8080"
	fmt.Printf("Servidor rodando em http://localhost:%s/\n", port)

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
