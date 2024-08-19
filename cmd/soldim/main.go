package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
    http.HandleFunc("/test-db-connection", func(w http.ResponseWriter, r *http.Request) {
        connStr := "host=c3gtj1dt5vh48j.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com port=5432 user=u8v8gfju2nbfqn password=p009827e7e08f28b44b9ba56751f82543345f445a718fb44624dd5b653e0238bd dbname=ddcvr5rele9132 sslmode=require"

        log.Println("About to connect to the database with the following connection string")
        log.Println("Connection string:", connStr)

        db, err := sql.Open("postgres", connStr)
        if err != nil {
            log.Println("Failed to open the database:", err)
            http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
            return
        }

        err = db.Ping()
        if err != nil {
            log.Println("Failed to connect to the database:", err)
            http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
            return
        }

        log.Println("Connected to the database successfully!")
        fmt.Fprintln(w, "Connected to the database successfully!")
    })

    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("$PORT must be set")
    }

    log.Printf("Server started on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
