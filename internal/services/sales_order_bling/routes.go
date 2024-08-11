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
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
)

const (
	limitePorPagina = 100
	bearerToken     = "a635e423f261781f4d0b8dcc710de4d6caa60d44"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_sales_orders", handleImportBlingSalesOrdersToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_sales_order_bling", handleGetSalesOrder).Methods(http.MethodGet)
}

func handleImportBlingSalesOrdersToSoldim(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			http.Error(w, fmt.Sprintf("Internal server error: %v", r), http.StatusInternalServerError)
		}
	}()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = limitePorPagina
	}

	fmt.Printf("Requesting page: %d with limit: %d\n", page, limit)

	for {
		sales, totalPages, err := bling.GetSalesOrdersFromBling(bearerToken, page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, sale := range sales {
			fmt.Printf("Verificando o contato com ID %d\n", sale.Contato.ID)
			contact, err := existContact(sale.Contato.ID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error checking contact existence: %v", err), http.StatusInternalServerError)
				return
			}

			if contact == nil || contact.ID == 0 {
				fmt.Printf("Contato com ID %d não encontrado. Criando novo contato.\n", sale.Contato.ID)
				newContact := &types.Contact{
					ID:              sale.Contato.ID,
					Nome:            sale.Contato.Nome,
					Codigo:          "",
					Situacao:        "",
					Numerodocumento: sale.Contato.NumeroDocumento,
					Telefone:        "",
					Celular:         "",
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				createdContact, err := createContact(*newContact)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error creating contact: %v", err), http.StatusInternalServerError)
					return
				}
				fmt.Printf("Contato criado com ID %d\n", createdContact.ID)
				sales[i].Contato.ID = createdContact.ID
				sales[i].ContactID = createdContact.ID
			} else {
				fmt.Printf("Contato encontrado com ID %d\n", contact.ID)
				sales[i].Contato.ID = contact.ID
				sales[i].ContactID = contact.ID
			}
		}

		fmt.Printf("Processing page: %d with %d products\n", page, len(sales))
		processSales(sales)

		if page >= totalPages {
			break
		}

		page++
	}

	err = processItemsSalesOrder(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar itens dos pedidos de venda: %v", err), http.StatusInternalServerError)
		return
	}

	err = updateSalesOrder()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func existContact(contactID int64) (*types.Contact, error) {
	url := fmt.Sprintf("http://localhost:8080/get_contact/%d", contactID)
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
	url := "http://localhost:8080/create_contact"
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

func processSales(sales []types.SalesOrder) {
	for _, sale := range sales {
		sale.SituationID = sale.Situacao.ID
		sale.StoreID = sale.Loja.ID
		salesOrderJSON, err := json.Marshal(sale)
		if err != nil {
			fmt.Printf("Error marshalling sales: %v\n", err)
			continue
		}

		// Adicione um log para imprimir o JSON que está sendo enviado
		fmt.Printf("SalesOrder JSON: %s\n", string(salesOrderJSON))
		req, err := http.NewRequest("POST", "http://localhost:8080/create_sales_order", bytes.NewBuffer(salesOrderJSON))
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

func processItemsSalesOrder(bearerToken string) error {
	// 1. Fazer a requisição para obter os IDs dos pedidos de venda
	resp, err := http.Get("http://localhost:8080/get_sales_orders")
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
		// Aqui fazemos a chamada para obter os detalhes do pedido de venda usando o ID
		salesOrderDetails, err := bling.GetSalesOrdersIDInBling(bearerToken, order.ID)
		if err != nil {
			fmt.Printf("Erro ao obter detalhes do pedido de venda com ID %d: %v\n", order.ID, err)
			continue
		}

		// 5. Processar cada item no pedido de venda retornado
		for _, item := range salesOrderDetails.Itens {
			// productID := item.ProductID.ID
			item.SalesOrderID = order.ID
			err := sendItemToCreate(item)
			if err != nil {
				fmt.Printf("Erro ao criar item %d do pedido %d: %v\n", item.ID, salesOrderDetails.ID, err)
				continue
			}
		}
	}

	return nil
}

func sendItemToCreate(item types.ItemsSalesOrders) error {
	orderItem := convertToItemsSalesOrder(item)

	url := "http://localhost:8080/create_items_sales_order"

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

	resp, err := http.Get("http://localhost:8080/get_sales_orders")
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

	resp, err := http.Get("http://localhost:8080/get_sales_orders")
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

	req, err := http.NewRequest("POST", "http://localhost:8080/create_products_sales_order", bytes.NewBuffer(productsalesorderJSON))
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

	req, err := http.NewRequest("POST", "http://localhost:8080/create_products_sales_order", bytes.NewBuffer(productsalesorderJSON))
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
	salesOrder, err := bling.GetSalesOrdersIDInBling(bearerToken, int64(salesOrderID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, salesOrder)
}
