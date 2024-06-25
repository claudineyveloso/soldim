package salechannel

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_sales_channels", handleGetSaleChannel).Methods(http.MethodGet)
}

func handleGetSaleChannel(w http.ResponseWriter, r *http.Request) {
	bearerToken := "5904172ee8be9f4f4ee15910b49752ca19c4b04c"
	channels, err := bling.GetSalesChannelsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter canais de vendas: %v", err), http.StatusInternalServerError)
		return
	}

	if len(channels) == 0 {
		http.Error(w, "Nenhum canal de vendas encontrado.", http.StatusNotFound)
		return
	}

	// Definir o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converter os canais de vendas para JSON e enviar a resposta
	err = json.NewEncoder(w).Encode(channels)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter canais de vendas para JSON: %v", err), http.StatusInternalServerError)
	}
}
