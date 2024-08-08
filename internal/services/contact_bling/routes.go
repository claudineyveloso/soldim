package contactbling

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
	bearerToken = "4e013e56e7ac5f1b915c3c68e3758c0624461a5f"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_contacts", handleImportBlingContactsToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_contacts_bling", handleGetContacts).Methods(http.MethodGet)
}

func handleImportBlingContactsToSoldim(w http.ResponseWriter, r *http.Request) {
	// Obtenha os contatos do Bling
	channels, err := bling.GetContactsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting contacts from Bling: %v", err), http.StatusInternalServerError)
		return
	}

	// Para cada contato, faça uma requisição para criar o contato no sistema local
	for _, channel := range channels {
		channelJSON, err := json.Marshal(channel)
		if err != nil {
			fmt.Printf("Error marshalling contacts: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/create_contact", bytes.NewBuffer(channelJSON))
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
			fmt.Printf("Failed to create contact. Status: %v, Response: %s\n", resp.Status, string(body))
			continue
		}

		fmt.Printf("Contact created successfully: %v\n", channel)
	}

	response := map[string]interface{}{
		"message": "Contatos importados e cadastrados com sucesso",
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

func handleGetContacts(w http.ResponseWriter, r *http.Request) {
	channels, err := bling.GetContactsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter contatos: %v", err), http.StatusInternalServerError)
		return
	}

	if len(channels) == 0 {
		http.Error(w, "Nenhum contato encontrado.", http.StatusNotFound)
		return
	}

	// Definir o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converter os depositos para JSON e enviar a resposta
	err = json.NewEncoder(w).Encode(channels)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter contatos para JSON: %v", err), http.StatusInternalServerError)
	}
}
