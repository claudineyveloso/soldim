package salesorderbling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/errors"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
)

const (
	limitePorPagina = 100
	bearerToken     = "4e013e56e7ac5f1b915c3c68e3758c0624461a5f"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_sales_orders", handleImportBlingSalesOrdersToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_sales_order_bling", handleGetSalesOrder).Methods(http.MethodGet)
}

func handleImportBlingSalesOrdersToSoldim(w http.ResponseWriter, r *http.Request) {
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
			contact, err := existContact(sale.Contato.ID)
			if err != nil {
				if err == errors.ErrContactNotFound {
					newContact := &types.Contact{
						ID:              sale.Contato.ID,
						Nome:            sale.Contato.Nome,
						Codigo:          "", // Adicione o código se disponível
						Situacao:        "", // Adicione a situação se disponível
						Numerodocumento: sale.Contato.NumeroDocumento,
						Telefone:        "", // Adicione o telefone se disponível
						Celular:         "", // Adicione o celular se disponível
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}
					createdContact, err := createContact(*newContact)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// Atualize o contact ID na venda com o ID retornado após a criação
					sales[i].Contato.ID = createdContact.ID
					sales[i].ContactID = createdContact.ID

				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				// Atualize o contact ID na venda com o ID encontrado
				sales[i].Contato.ID = contact.ID
				sales[i].ContactID = contact.ID
			}
		}

		fmt.Printf("Processing page: %d with %d products\n", page, len(sales))
		processSales(sales) // Processa todas as vendas

		if page >= totalPages {
			break
		}

		page++
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

	// Decodifica a resposta para verificar se há uma mensagem de erro
	var errorResponse ErrorResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Verifica se a resposta contém um erro
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorResponse.Error == "contact not found" {
			fmt.Printf("Contato com ID %d não encontrado.\n", contactID)
			return nil, errors.ErrContactNotFound
		}
	}

	// Verifica se a resposta contém o contato
	var contact types.Contact
	if err := json.Unmarshal(body, &contact); err != nil {
		return nil, fmt.Errorf("error decoding contact response: %v", err)
	}

	// Log detalhado do contato encontrado
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
