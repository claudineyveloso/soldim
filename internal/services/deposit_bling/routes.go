package depositbling

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
	bearerToken = "3a1a2228e56718a1e220428cda2b7bce2454281b"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_deposits", handleImportBlingDepositsToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_deposits_bling", handleGetDeposits).Methods(http.MethodGet)
}

func handleImportBlingDepositsToSoldim(w http.ResponseWriter, r *http.Request) {
	// Obtenha os depositos do Bling
	channels, err := bling.GetDepositsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting deposits from Bling: %v", err), http.StatusInternalServerError)
		return
	}

	// Para cada deposito, faça uma requisição para criar o deposito no sistema local
	for _, channel := range channels {
		channelJSON, err := json.Marshal(channel)
		if err != nil {
			fmt.Printf("Error marshalling deposits: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/create_deposit", bytes.NewBuffer(channelJSON))
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
			fmt.Printf("Failed to create deposit. Status: %v, Response: %s\n", resp.Status, string(body))
			continue
		}

		fmt.Printf("Deposit created successfully: %v\n", channel)
	}

	response := map[string]interface{}{
		"message": "Depositos importados e cadastrados com sucesso",
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

func handleGetDeposits(w http.ResponseWriter, r *http.Request) {
	channels, err := bling.GetDepositsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter depositos: %v", err), http.StatusInternalServerError)
		return
	}

	if len(channels) == 0 {
		http.Error(w, "Nenhum deposito encontrado.", http.StatusNotFound)
		return
	}

	// Definir o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converter os depositos para JSON e enviar a resposta
	err = json.NewEncoder(w).Encode(channels)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter depositos para JSON: %v", err), http.StatusInternalServerError)
	}
}
