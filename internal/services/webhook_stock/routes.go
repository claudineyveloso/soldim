package webhookstock

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// Registra o endpoint do webhook
	router.HandleFunc("/webhook_stock", handleBlingWebhookStock).Methods(http.MethodPost)
}

func handleBlingWebhookStock(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Loga o recebimento da requisição
	log.Println("Recebendo webhook de atualização de produto da Bling...")

	// Lê o corpo da requisição
	if err := r.ParseForm(); err != nil {
		log.Printf("Erro ao analisar o corpo da requisição: %v", err)
		http.Error(w, "Erro ao analisar o corpo da requisição", http.StatusInternalServerError)
		return
	}

	// Obtém o valor do parâmetro 'data', que contém o JSON
	encodedData := r.FormValue("data")
	if encodedData == "" {
		log.Println("Parâmetro 'data' não encontrado no corpo da requisição")
		http.Error(w, "Parâmetro 'data' ausente", http.StatusBadRequest)
		return
	}

	// Decodifica a string 'urlencoded' para JSON
	decodedData, err := url.QueryUnescape(encodedData)
	if err != nil {
		log.Printf("Erro ao decodificar a string URL: %v", err)
		http.Error(w, "Erro ao decodificar o parâmetro 'data'", http.StatusBadRequest)
		return
	}

	// Loga o JSON bruto decodificado
	log.Printf("JSON decodificado recebido: %s\n", decodedData)

	// Decodifica o JSON para a struct BlingWebhookResponse
	var response types.BlingWebhookResponse
	if err := json.Unmarshal([]byte(decodedData), &response); err != nil {
		log.Printf("Erro ao decodificar JSON do webhook: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Loga o conteúdo do produto para depuração
	log.Printf("Dados do estoque recebidos: %+v\n", response)

	// Itera sobre os estoques recebidos
	for _, estoqueWrapper := range response.Retorno.Estoques {
		estoque := estoqueWrapper.Estoque
		log.Printf("Atualizando estoque ID: %d, Código: %s, Nome: %s, Estoque Atual: %d\n",
			estoque.ID, estoque.Codigo, estoque.Nome, estoque.EstoqueAtual)

		// Chama a função para atualizar o estoque na base de dados
		if err := updateProductInDatabase(estoque); err != nil {
			log.Printf("Falha ao atualizar o produto na base de dados: %v", err)
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}
	}

	// Retorna sucesso
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func updateProductInDatabase(estoque types.Estoque) error {
	// Lógica para atualizar a base de dados no Soldim
	// Pode ser uma chamada ao banco de dados ou uma API

	// Log fictício para mostrar a intenção
	log.Printf("Atualizando o produto na base de dados: %+v\n", estoque)

	return nil
}
