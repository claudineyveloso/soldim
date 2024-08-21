package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// handler para a rota "/"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Obtém a porta do ambiente do Heroku
	port := os.Getenv("PORT")
	if port == "" {
		// Se não estiver definida, usa a porta 8080 por padrão
		port = "8080"
	}

	http.HandleFunc("/", helloHandler)

	fmt.Printf("Servidor rodando em http://localhost:%s/\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
