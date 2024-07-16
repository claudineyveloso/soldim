package productbling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

const (
	limitePorPagina = 100
	bearerToken     = "b925150715499ff7c242e9f0b298d92c69ab0bf3"
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

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			processStocks(products, bearerToken)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			processSuppliers(products, bearerToken)
		}()

		wg.Wait()

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

		fmt.Println("***********************************************************************************")
		fmt.Printf("Updating product with ID: %d\n", blingProduct.ID)
		fmt.Println("***********************************************************************************")

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

		// Adicionar log detalhado do updateProduct

		fmt.Println("***********************************************************************************")
		fmt.Printf("Update product details: %+v\n", updateProduct)
		fmt.Println("***********************************************************************************")

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
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
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

func processStocks(products []types.Product, bearerToken string) {
	var wg sync.WaitGroup

	// Criar um rate limiter que permite 3 requisições por segundo
	limiter := rate.NewLimiter(rate.Every(time.Second/3), 1)

	for _, product := range products {
		wg.Add(1)
		go func(product types.Product) {
			defer wg.Done()

			// Aguardar até que uma requisição possa ser feita
			err := limiter.Wait(context.Background())
			if err != nil {
				fmt.Printf("Error waiting for rate limiter: %v\n", err)
				return
			}

			processStockForProduct(product, bearerToken)
		}(product)
	}

	wg.Wait()
}

func processStockForProduct(product types.Product, bearerToken string) {
	stockResponse, err := bling.GetStockProductFromBling(bearerToken, product.ID)
	if err != nil {
		fmt.Printf("Error fetching stock for product %d: %v\n", product.ID, err)
		utils.LogErrorToFile(fmt.Sprintf("Error fetching stock for product %d: %v\n", product.ID, err))

		return
	}

	// fmt.Printf("Parsed stock response for product %d: %+v\n", product.ID, stockResponse)

	for _, stockData := range stockResponse.Data {
		// fmt.Printf("Processing stock data: %+v\n", stockData)

		// Criar o stock
		stock := types.Stock{
			ProductID:         stockData.Produto.ID,
			SaldoFisicoTotal:  int32(stockData.SaldoFisicoTotal),
			SaldoVirtualTotal: int32(stockData.SaldoVirtualTotal),
		}

		stockJSON, err := json.Marshal(stock)
		if err != nil {
			fmt.Printf("Error marshalling stock for product %d: %v\n", product.ID, err)
			continue
		}

		// fmt.Printf("Sending stock data for product %d: %s\n", product.ID, string(stockJSON))

		stockResp, err := http.Post("http://localhost:8080/create_stock", "application/json", bytes.NewBuffer(stockJSON))
		if err != nil {
			fmt.Printf("Error sending stock data for product %d: %v\n", product.ID, err)
			continue
		}
		defer stockResp.Body.Close()

		// Criar deposit products
		for _, deposito := range stockData.Depositos {
			// fmt.Printf("Processing deposit data: %+v\n", deposito)

			depositProduct := types.DepositProduct{
				ProductID:    stockData.Produto.ID,
				DepositID:    deposito.ID,
				SaldoFisico:  int32(deposito.SaldoFisico),
				SaldoVirtual: int32(deposito.SaldoVirtual),
			}

			// fmt.Printf("Values assigned for deposit product for product %d: %+v\n", product.ID, depositProduct)

			depositProductJSON, err := json.Marshal(depositProduct)
			if err != nil {
				fmt.Printf("Error marshalling deposit product for product %d: %v\n", product.ID, err)
				continue
			}

			// fmt.Printf("Sending deposit product data for product %d: %s\n", product.ID, string(depositProductJSON))

			depositResp, err := http.Post("http://localhost:8080/create_deposit_product", "application/json", bytes.NewBuffer(depositProductJSON))
			if err != nil {
				fmt.Printf("Error sending deposit product data for product %d: %v\n", product.ID, err)
				continue
			}
			defer depositResp.Body.Close()

			// depositRespBody, _ := io.ReadAll(depositResp.Body)
			// fmt.Printf("Response from create_deposit_product for product %d: %s\n", product.ID, string(depositRespBody))
		}
	}
}

func processSuppliers(products []types.Product, bearerToken string) {
	var wg sync.WaitGroup

	// Criar um rate limiter que permite 3 requisições por segundo
	limiter := rate.NewLimiter(rate.Every(time.Second/3), 1)

	for _, product := range products {
		wg.Add(1)
		go func(product types.Product) {
			defer wg.Done()

			// Aguardar até que uma requisição possa ser feita
			err := limiter.Wait(context.Background())
			if err != nil {
				fmt.Printf("Error waiting for rate limiter: %v\n", err)
				return
			}

			processSupplierForProduct(product, bearerToken)
		}(product)
	}

	wg.Wait()
}

func processSupplierForProduct(product types.Product, bearerToken string) {
	supplierResponse, err := bling.GetSupplierProductFromBling(bearerToken, product.ID)
	if err != nil {
		fmt.Printf("Error fetching supplier for product %d: %v\n", product.ID, err)
		return
	}

	for _, supplierData := range supplierResponse.Data {
		// Criar o supplier product
		supplierProduct := types.SupplierProduct{
			ID:          supplierData.ID,
			Descricao:   supplierData.Descricao,
			PrecoCusto:  supplierData.PrecoCusto,
			PrecoCompra: supplierData.PrecoCompra,
			Padrao:      supplierData.Padrao,
			SupplierID:  supplierData.Fornecedor.ID,
			ProductID:   supplierData.Produto.ID,
		}

		supplierProductJSON, err := json.Marshal(supplierProduct)
		if err != nil {
			fmt.Printf("Error marshalling supplier product for product %d: %v\n", product.ID, err)
			continue
		}

		fmt.Printf("Sending supplier product data for product %d: %s\n", product.ID, string(supplierProductJSON))

		supplierProductResp, err := http.Post("http://localhost:8080/create_supplier_product", "application/json", bytes.NewBuffer(supplierProductJSON))
		if err != nil {
			fmt.Printf("Error sending supplier product data for product %d: %v\n", product.ID, err)
			continue
		}
		defer supplierProductResp.Body.Close()

		supplierProductRespBody, _ := io.ReadAll(supplierProductResp.Body)
		fmt.Printf("Response from create_supplier_product for product %d: %s\n", product.ID, string(supplierProductRespBody))
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
