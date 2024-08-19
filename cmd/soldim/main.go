package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
    connStr := "host=c3gtj1dt5vh48j.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com port=5432 user=u8v8gfju2nbfqn password=p009827e7e08f28b44b9ba56751f82543345f445a718fb44624dd5b653e0238bd dbname=ddcvr5rele9132 sslmode=require"

    log.Println("About to connect to the database with the following connection string")
    log.Println("Connection string:", connStr)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    log.Println("Connected to the database successfully!")
}
