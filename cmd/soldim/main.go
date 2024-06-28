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
	cfg := configs.Config{
		PublicHost: configs.Envs.PublicHost,
		Port:       configs.Envs.Port,
		DBUser:     configs.Envs.DBUser,
		DBPassword: configs.Envs.DBPassword,
		DBName:     configs.Envs.DBName,
	}
	db, err := db.NewPostgresSQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")
}

// package main
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
//
// 	"github.com/claudineyveloso/soldim.git/internal/crawler"
// 	"github.com/gorilla/mux"
// )
//
// func main() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/crawl", handleCrawl).Methods("GET")
//
// 	fmt.Println("Server started at :8080")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }
//
// func handleCrawl(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("query")
// 	if query == "" {
// 		http.Error(w, "query parameter is required", http.StatusBadRequest)
// 		return
// 	}
//
// 	produtos, err := crawler.CrawlGoogle(query)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error crawling data: %v", err), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(produtos)
//
// 	// Log the total number of products collected
// 	log.Printf("Total de produtos coletados: %d", len(produtos))
// }
