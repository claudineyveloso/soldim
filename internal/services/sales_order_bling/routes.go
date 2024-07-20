package salesorderbling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/gorilla/mux"
)

const (
	limitePorPagina = 100
	bearerToken     = "9886c96e963dac41b5c3b24c12dd03751ce9d9f7"
)

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

	fmt.Printf("Requesting page: %d with limit: %d", page, limit)

	for {
		sales, totalPages, err := bling.GetSalesOrdersFromBling(bearerToken, page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Processing page: %d with %d products\n", page, len(sales))
		processSales(sales)

		if page >= totalPages {
			break
		}

		page++

	}
	//err = processProductsSalesOrders()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
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
			fmt.Printf("Failed to create sales orders. Status: %v\n", resp.Status)
			continue
		}

		fmt.Printf("Sales Orders created successfully: %v\n", sale)
	}
}

// func processProductsSalesOrders() error {
// 	resp, err := http.Get("http://localhost:8080/get_sales_orders")
// 	if err != nil {
// 		return fmt.Errorf("erro ao chamar get_sales_orders: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("falha na requisição para get_sales_orders: %s", resp.Status)
// 	}
//
// 	var salesOrders []types.SalesOrder
// 	err = json.NewDecoder(resp.Body).Decode(&salesOrders)
// 	if err != nil {
// 		return fmt.Errorf("erro ao decodificar resposta: %v", err)
// 	}
//
// 	for _, salesOrder := range salesOrders {
// 		salesOrderData, err := bling.GetSalesOrdersIDInBling(bearerToken, salesOrder.ID)
// 		if err != nil {
// 			return fmt.Errorf("erro ao obter SalesOrdersIDInBling para ID %d: %v", salesOrder.ID, err)
// 		}
//
// 		err = createProductsSalesOrder(salesOrderData)
// 		if err != nil {
// 			return fmt.Errorf("erro ao criar ProductSalesOrder para ID %d: %v", salesOrder.ID, err)
// 		}
// 	}
//
// 	return nil
// }

// func createProductsSalesOrder(salesOrderData *types.BlingSalesOrderData) error {
// 	for _, item := range salesOrderData.Itens {
// 		// Implemente a lógica de criação de ProductSalesOrder aqui, usando item.ID e item.Quantidade
// 		fmt.Printf("Criando ProductSalesOrder para Item ID: %d, Quantidade: %d\n", item.ID, item.Quantidade)
// 	}
// 	return nil
// }

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
