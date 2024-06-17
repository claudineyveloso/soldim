package healthy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/healthy", handleGetHealthy).Methods(http.MethodGet)
}

func handleGetHealthy(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Bem vindo ao Soldim!"))
	if err != nil {
		http.Error(w, "Erro ao escrever resposta", http.StatusInternalServerError)
		return
	}
}
