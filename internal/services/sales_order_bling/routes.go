package salesorderbling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/pkg/utils"
	"github.com/gorilla/mux"
)

var (
	limitePorPagina = 100
	baseURL         = utils.GetBaseURL()
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_sales_orders", handleImportBlingSalesOrdersToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_sales_order_bling", handleGetSalesOrder).Methods(http.MethodGet)
}

func handleImportBlingSalesOrdersToSoldim(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 100

	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar token de acesso: %v", err), http.StatusInternalServerError)
		return
	}
	rateLimiter := time.NewTicker(333 * time.Millisecond)
	defer rateLimiter.Stop()

	fmt.Printf("Requesting page: %d with limit: %d\n", page, limit)

	for {
		<-rateLimiter.C

		sales, totalPages, err := bling.GetSalesOrdersFromBling(token, page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(sales) == 0 {
			break
		}
		processSalesOrdersConcurrently(sales, rateLimiter, token)

		if page >= totalPages {
			break
		}

		page++
	}
	responseMessage := map[string]interface{}{
		"message": "Registros importados e atualizados com sucesso",
		"status":  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(responseMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func processSalesOrdersConcurrently(salesorders []types.SalesOrder, rateLimiter *time.Ticker, token string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // Limita a 3 goroutines simultâneas
	for _, sale := range salesorders {
		wg.Add(1)
		saleCopy := sale // Criação de uma cópia para evitar problemas com concorrência
		go func(sale types.SalesOrder) {
			defer wg.Done()
			sem <- struct{}{} // Adquire um recurso do semáforo para limitar as goroutines
			processSalesOrders(sale, rateLimiter, token)
			<-sem // Libera o recurso do semáforo
		}(saleCopy) // Passa a cópia da variável `sale` para a goroutine
	}
	wg.Wait() // Aguarda todas as goroutines terminarem
}

func processSalesOrders(sale types.SalesOrder, rateLimiter *time.Ticker, token string) {
	const maxRetries = 5
	backoff := time.Second
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		<-rateLimiter.C // Espera o próximo "tick" do rateLimiter antes de processar
		salesOrder, err := bling.GetSalesOrdersIDInBling(token, sale.ID)
		if err != nil {
			if isRateLimitError(err) { // Verifica se o erro é devido ao limite de requisições
				fmt.Printf("Limite de requisições atingido, esperando antes de tentar novamente...\n")
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			fmt.Printf("Erro ao obter detalhes do pedido de venda com ID %d: %v\n", sale.ID, err)
			return
		}
		utils.ProcessContact(*salesOrder)
		contactID := salesOrder.Contato.ID
		utils.ProcessSales(*salesOrder, contactID)
		utils.ProcessAllItems(*salesOrder)
		return
	}
}

func isRateLimitError(err error) bool {
	// Implementar verificação de erro de limite se o Bling retornar um código de status específico ou mensagem
	// Exemplo: return strings.Contains(err.Error(), "rate limit")
	if strings.Contains(err.Error(), "rate limit") {
		return true
	}
	return false
}

func existContact(contactID int64) (*types.Contact, error) {
	url := fmt.Sprintf(baseURL+"/get_contact/%d", contactID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching contact: %v", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta antes de ler o corpo
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Contato com ID %d não encontrado.\n", contactID)
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Leitura do corpo para fins de depuração, mas ignorada
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Ler o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Verifica se o corpo da resposta é "null" ou "{}"
	bodyTrimmed := strings.TrimSpace(string(body))
	if bodyTrimmed == "null" || bodyTrimmed == "{}" {
		fmt.Printf("Contato com ID %d não encontrado.\n", contactID)
		return nil, nil
	}

	// Tenta decodificar o corpo da resposta como um contato
	var contact types.Contact
	if err := json.Unmarshal(body, &contact); err != nil {
		return nil, fmt.Errorf("error decoding contact response: %v", err)
	}

	fmt.Printf("Contato encontrado: %+v\n", contact)
	return &contact, nil
}

func createContact(contact types.Contact) (*types.Contact, error) {
	url := baseURL + "/create_contact"
	contactData, err := json.Marshal(contact)
	if err != nil {
		return nil, fmt.Errorf("error marshaling contact data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(contactData))
	if err != nil {
		return nil, fmt.Errorf("error creating contact: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var createdContact types.Contact
	err = json.NewDecoder(resp.Body).Decode(&createdContact)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	// Verifica se o contato realmente existe no banco de dados após a criação
	existingContact, err := existContact(createdContact.ID)
	if err != nil {
		return nil, fmt.Errorf("contact not found after creation: %v", err)
	}

	return existingContact, nil
}

func processSales(sales []types.SalesOrder, rateLimiter *time.Ticker) {
	for _, sale := range sales {
		<-rateLimiter.C
		sale.SituationID = sale.Situacao.ID
		sale.StoreID = sale.Loja.ID
		salesOrderJSON, err := json.Marshal(sale)
		if err != nil {
			fmt.Printf("Error marshalling sales: %v\n", err)
			continue
		}

		// Adicione um log para imprimir o JSON que está sendo enviado
		fmt.Printf("SalesOrder JSON: %s\n", string(salesOrderJSON))
		req, err := http.NewRequest("POST", baseURL+"/create_sales_order", bytes.NewBuffer(salesOrderJSON))
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
			fmt.Printf("Failed to create sales orders. Status: %v, Response: %s\n", resp.Status, string(body))
			continue
		}

		fmt.Printf("Sales Orders created successfully em processSales: %v\n", sale)
	}
}

func processItemsSalesOrder(bearerToken string, rateLimiter *time.Ticker) error {
	// 1. Fazer a requisição para obter os IDs dos pedidos de venda
	resp, err := http.Get(baseURL + "/get_sales_orders")
	if err != nil {
		return fmt.Errorf("erro ao chamar get_sales_orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha na requisição para get_sales_orders: %s", resp.Status)
	}

	// 2. Ler a resposta do corpo
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	// 3. Desserializar o JSON na estrutura correta
	var response types.SalesOrdersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("erro ao desserializar a resposta: %v", err)
	}

	// 4. Iterar sobre cada pedido de venda para processar os itens
	for _, order := range response.SalesOrders {
		<-rateLimiter.C
		// Aqui fazemos a chamada para obter os detalhes do pedido de venda usando o ID
		salesOrderDetails, err := bling.GetSalesOrdersIDInBling(bearerToken, order.ID)
		if err != nil {
			fmt.Printf("Erro ao obter detalhes do pedido de venda com ID %d: %v\n", order.ID, err)
			continue
		}

		// 5. Processar cada item no pedido de venda retornado
		for _, item := range salesOrderDetails.Itens {
			<-rateLimiter.C
			fmt.Printf("Processando item: %+v\n", item)
			// productID := item.ProductID.ID
			item.SalesOrderID = order.ID
			err := sendItemToCreate(item)
			if err != nil {
				fmt.Printf("Erro ao criar item %d do pedido %d: %v\n", item.ID, salesOrderDetails.ID, err)
				continue
			}
			fmt.Printf("Item %d do pedido %d criado com sucesso.\n", item.ID, salesOrderDetails.ID)
		}
	}

	return nil
}

func sendItemToCreate(item types.ItemsSalesOrders) error {
	orderItem := convertToItemsSalesOrder(item)

	url := baseURL + "/create_items_sales_order"

	// Cria o payload a partir do item, que agora inclui o product_id corretamente
	payload, err := json.Marshal(orderItem)
	if err != nil {
		return fmt.Errorf("erro ao serializar o item: %v", err)
	}

	// Cria uma requisição POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição POST: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("falha ao criar item: %s, resposta: %s", resp.Status, string(body))
	}

	return nil
}

func convertToItemsSalesOrder(item types.ItemsSalesOrders) types.ItemsSalesOrder {
	return types.ItemsSalesOrder{
		ID:                 item.ID,
		SalesOrderID:       item.SalesOrderID,
		Codigo:             item.Codigo,
		Unidade:            item.Unidade,
		Quantidade:         item.Quantidade,
		Desconto:           item.Desconto,
		Valor:              item.Valor,
		Aliquotaipi:        item.Aliquotaipi,
		Descricao:          item.Descricao,
		Descricaodetalhada: item.Descricaodetalhada,
		ProductID:          item.ProductID.ID, // Extraia o ID do produto corretamente
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}
}

func updateSalesOrder() error {
	logFile, err := os.OpenFile("error_import_sales_orders_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}
	defer logFile.Close()

	resp, err := http.Get(baseURL + "/get_sales_orders")
	if err != nil {
		return fmt.Errorf("erro ao chamar get_sales_orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha na requisição para get_sales_orders: %s", resp.Status)
	}
	return nil
}

func processProductsSalesOrders() error {
	logFile, err := os.OpenFile("error_import_sales_orders_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}
	defer logFile.Close()

	resp, err := http.Get(baseURL + "/get_sales_orders")
	if err != nil {
		return fmt.Errorf("erro ao chamar get_sales_orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha na requisição para get_sales_orders: %s", resp.Status)
	}

	// for _, salesOrder := range salesOrders {
	// 	salesOrderData, err := bling.GetSalesOrdersIDInBling(bearerToken, salesOrder.ID)
	// 	if err != nil {
	// 		// return fmt.Errorf("erro ao obter SalesOrdersIDInBling para ID %d: %v", salesOrder.ID, err)
	// 		utils.LogError(logFile, salesOrder.ID, fmt.Errorf("erro ao obter SalesOrdersIDInBling para ID %d: %v", salesOrder.ID, err))
	// 		continue
	// 	}
	// 	fmt.Printf("Na linha abaixo, pegar todos os pedidos de venda dentro do FOR")
	//
	// 	for _, item := range salesOrderData.Itens {
	// 		productSalesOrder := types.ProductSalesOrderPayload{
	// 			SalesOrderID: salesOrderData.ID,
	// 			ProductID:    item.Produto.ID,
	// 			Quantidade:   int32(item.Quantidade), // Converte quantidade para int32
	// 			CreatedAt:    time.Now(),
	// 			UpdatedAt:    time.Now(),
	// 		}
	//
	// 		err = createProductsSalesOrder(productSalesOrder)
	// 		if err != nil {
	// 			utils.LogError(logFile, salesOrder.ID, fmt.Errorf("erro ao criar ProductSalesOrder para SalesOrder ID %d e Item ID %d: %v", salesOrder.ID, item.ID, err))
	// 			continue
	// 			// return fmt.Errorf("erro ao criar ProductSalesOrder para SalesOrder ID %d e Item ID %d: %v", salesOrder.ID, item.ID, err)
	// 		}
	// 	}
	// }

	return nil
}

func createProductsSalesOrder(productsalesorder types.ProductSalesOrderPayload) error {
	productsalesorderJSON, err := json.Marshal(productsalesorder)
	if err != nil {
		return fmt.Errorf("error marshalling product sales order: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/create_products_sales_order", bytes.NewBuffer(productsalesorderJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create product sales orders. Status: %v, Response: %s", resp.Status, string(body))
	}

	fmt.Printf("Product Sales Orders created successfully: %v\n", productsalesorder)
	fmt.Printf("Response Body: %s\n", string(body))
	return nil
}

func createProductsSalesOrder_OLD(productsalesorder types.ProductSalesOrderPayload) error {
	productsalesorderJSON, err := json.Marshal(productsalesorder)
	if err != nil {
		return fmt.Errorf("error marshalling product sales order: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/create_products_sales_order", bytes.NewBuffer(productsalesorderJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create product sales orders. Status: %v", resp.Status)
	}

	fmt.Printf("Product Sales Orders created successfully: %v\n", productsalesorder)
	return nil
}

func handleGetSalesOrder(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar token de acesso: %v", err), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	salesOrderIDStr, ok := vars["salesOrderID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Produto ausente!"))
		return
	}
	salesOrderID, err := strconv.Atoi(salesOrderIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Pedido de Vendas inválido: %v", err))
		return
	}
	salesOrder, err := bling.GetSalesOrdersIDInBling(token, int64(salesOrderID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, salesOrder)
}
