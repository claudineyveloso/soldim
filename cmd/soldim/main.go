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
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
