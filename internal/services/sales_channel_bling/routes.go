package saleschannelbling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
)

const (
	bearerToken = "6ba9cd2953cee77ad3a88b354cbb67ef72c1012b"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_sales_channel", handleImportBlingSalesChannelToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_sales_channel_bling", handleGetSaleChannel).Methods(http.MethodGet)
}

func handleImportBlingSalesChannelToSoldim(w http.ResponseWriter, r *http.Request) {
	// Obtenha os canais de venda do Bling
	channels, err := bling.GetSalesChannelsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting sales channels from Bling: %v", err), http.StatusInternalServerError)
		return
	}

	// Para cada canal de venda, faça uma requisição para criar o canal de venda no sistema local
	for _, channel := range channels {
		channelJSON, err := json.Marshal(channel)
		if err != nil {
			fmt.Printf("Error marshalling sales channel: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/create_sales_channel", bytes.NewBuffer(channelJSON))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Failed to create sales channel. Status: %v, Response: %s\n", resp.Status, string(body))
			continue
		}

		fmt.Printf("Sales channel created successfully: %v\n", channel)
	}

	response := map[string]interface{}{
		"message": "Canais de venda importados e cadastrados com sucesso",
		"status":  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func handleGetSaleChannel(w http.ResponseWriter, r *http.Request) {
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
