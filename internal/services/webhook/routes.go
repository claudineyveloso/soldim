package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	productbling "github.com/claudineyveloso/soldim.git/internal/services/product_bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/webhook/bling", handleBlingWebhook).Methods(http.MethodPost)
	productbling.RegisterRoutes(router)
}

func handleBlingWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookData struct {
		Produto struct {
			Codigo    string  `json:"codigo"`
			Descricao string  `json:"descricao"`
			Preco     float64 `json:"preco"`
			Tipo      string  `json:"tipo"`
		} `json:"produto"`
	}

	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao decodificar o webhook: %v", err), http.StatusBadRequest)
		return
	}

	// Criar um novo produto com os dados recebidos
	product := types.Product{
		Codigo: webhookData.Produto.Codigo,
		Nome:   webhookData.Produto.Descricao,
		Preco:  webhookData.Produto.Preco,
		Tipo:   webhookData.Produto.Tipo,
	}

	// Inserir o produto na base de dados local (comentado para fins de exemplo)
	// if err := db.Create(&product).Error; err != nil {
	//		http.Error(w, fmt.Sprintf("Erro ao salvar o produto: %v", err), http.StatusInternalServerError)
	//		return
	//	}

	// Log que um novo produto foi cadastrado
	log.Printf("Novo produto cadastrado: %s - %s - %.2f - %s",
		product.Codigo, product.Nome, product.Preco, product.Tipo)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Produto cadastrado com sucesso"))
}
