package productbling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
)

const (
	limitePorPagina = 100
	bearerToken     = "8d9851265ae849322827e320972b03c3596b69df" // r.Header.Get("Authorization")
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_products", handleImportBlingProductsToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_products_bling", handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/create_product_bling", handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/update_product_bling", handleUpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/delete_product_bling", handleDeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/get_product_id_bling", handleGetProductId).Methods(http.MethodGet)
}

func handleImportBlingProductsToSoldim(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	criterioStr := "5"
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = limitePorPagina
	}

	criterio, err := strconv.Atoi(criterioStr)
	if err != nil {
		criterio = 0 // Valor padrão para criterio se não for fornecido ou inválido
	}

	fmt.Printf("Requesting page: %d with limit: %d and name: %s\n", page, limit, name)

	for {
		products, totalPages, err := bling.GetProductsFromBling(bearerToken, page, limit, name, criterio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Processing page: %d with %d products\n", page, len(products))
		processProducts(products)

		if page >= totalPages {
			break
		}

		page++
	}

	resp, err := http.Get("http://localhost:8080/get_products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var products []types.Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, product := range products {
		// Substituímos a linha para chamar o endpoint local
		url := fmt.Sprintf("http://localhost:8080/get_product_id_bling?productID=%d", product.ID)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error getting product from Bling: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			continue
		}

		var blingProduct types.Product
		err = json.Unmarshal(body, &blingProduct)
		if err != nil {
			fmt.Printf("Error unmarshalling product: %v\n", err)
			continue
		}

		fmt.Printf("Updating product with ID: %d\n", blingProduct.ID)

		// Mapear BlingProduct para a estrutura necessária
		updateProduct := types.Product{
			ID:                         blingProduct.ID,
			Nome:                       blingProduct.Nome,
			Codigo:                     blingProduct.Codigo,
			Preco:                      blingProduct.Preco,
			Tipo:                       blingProduct.Tipo,
			Situacao:                   blingProduct.Situacao,
			Formato:                    blingProduct.Formato,
			DescricaoCurta:             blingProduct.DescricaoCurta,
			Datavalidade:               blingProduct.Datavalidade,
			Unidade:                    blingProduct.Unidade,
			Pesoliquido:                blingProduct.Pesoliquido,
			Pesobruto:                  blingProduct.Pesobruto,
			Volumes:                    blingProduct.Volumes,
			Itensporcaixa:              blingProduct.Itensporcaixa,
			Gtin:                       blingProduct.Gtin,
			Gtinembalagem:              blingProduct.Gtinembalagem,
			Tipoproducao:               blingProduct.Tipoproducao,
			Condicao:                   blingProduct.Condicao,
			Fretegratis:                blingProduct.Fretegratis,
			Marca:                      blingProduct.Marca,
			Descricaocomplementar:      blingProduct.Descricaocomplementar,
			Linkexterno:                blingProduct.Linkexterno,
			Observacoes:                blingProduct.Observacoes,
			Descricaoembalagemdiscreta: blingProduct.Descricaoembalagemdiscreta,
		}

		// Atualize o produto em localhost com os dados obtidos do Bling
		updateURL := fmt.Sprintf("http://localhost:8080/update_product?productID=%d", updateProduct.ID)
		productJSON, err := json.Marshal(updateProduct)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling updated product: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Product JSON: %s\n", string(productJSON)) // Adicionando log do JSON do produto

		req, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(productJSON))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating update request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending update request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Failed to update product. Status: %v, Response: %s\n", resp.Status, string(body)) // Adicionando log da resposta de erro
			http.Error(w, fmt.Sprintf("Failed to update product. Status: %v", resp.Status), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Product updated successfully: %v\n", updateProduct)
	}

	response := map[string]interface{}{
		"message": "Registros importados e atualizados com sucesso",
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

func processProducts(products []types.Product) {
	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			fmt.Printf("Error marshalling product: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/create_product", bytes.NewBuffer(productJSON))
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
			fmt.Printf("Failed to create product. Status: %v\n", resp.Status)
			continue
		}

		fmt.Printf("Product created successfully: %v\n", product)
	}
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	// bearerToken := "981b387171e4db2550a80c80eb1fbd7c6af0a807" // r.Header.Get("Authorization")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	criterioStr := r.URL.Query().Get("criterio")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = limitePorPagina
	}

	criterio, err := strconv.Atoi(criterioStr)
	if err != nil {
		criterio = 0 // Valor padrão para criterio se não for fornecido ou inválido
	}

	fmt.Printf("Requesting page: %d with limit: %d and name: %s\n", page, limit, name)
	products, totalPages, err := bling.GetProductsFromBling(bearerToken, page, limit, name, criterio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Products   []types.Product `json:"products"`
		TotalPages int             `json:"totalPages"`
		Page       int             `json:"page"`
		Limit      int             `json:"limit"`
	}{
		Products:   products,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retornando %d produtos e %d páginas\n", len(products), totalPages)
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o JSON do corpo da requisição para a estrutura Product
	var newProduct types.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Fecha o corpo da requisição após o processamento
	defer r.Body.Close()

	// bearerToken := "981b387171e4db2550a80c80eb1fbd7c6af0a807" // r.Header.Get("Authorization")
	// Chama a função para criar o produto no Bling
	err := bling.CreateProductInBling(bearerToken, newProduct)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao criar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Produto criado com sucesso!")
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da requisição para a estrutura Product
	var updatedProduct types.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Fecha o corpo da requisição após o processamento
	defer r.Body.Close()

	// bearerToken := "981b387171e4db2550a80c80eb1fbd7c6af0a807" // r.Header.Get("Authorization")
	// Chama a função para atualizar o produto no Bling
	err = bling.UpdateProductInBling(bearerToken, productID, updatedProduct)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao atualizar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produto atualizado com sucesso!")
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// bearerToken := "981b387171e4db2550a80c80eb1fbd7c6af0a807" // r.Header.Get("Authorization")
	// Chama a função para deletar o produto no Bling
	err = bling.DeleteProductInBling(bearerToken, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao deletar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao deletar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produto deletado com sucesso!")
}

func handleGetProductId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// bearerToken := "ea2648642cd55fa59ac6582d3a9506be8e91f6f2" // r.Header.Get("Authorization")
	// Chama a função para obter os detalhes do produto no Bling
	product, err := bling.GetProductIDInBling(bearerToken, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter detalhes do produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao obter detalhes do produto: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
