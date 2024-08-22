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
    "os"
)

func main() {
    // Use a porta fornecida pelo Heroku ou uma porta padrão (8080)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, world!")
    })

    log.Printf("Listening on port %s...", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatal(err)
    }
}
