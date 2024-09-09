package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func ExistContact(contactID int64) (*types.Contato, error) {
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
	var contato types.Contato
	if err := json.Unmarshal(body, &contato); err != nil {
		return nil, fmt.Errorf("error decoding contact response: %v", err)
	}

	fmt.Printf("Contato encontrado: %+v\n", contato)
	return &contato, nil
}

func ProcessSales(sales types.SalesOrder, contactID int64) {
	sales.SituationID = sales.Situacao.ID
	sales.StoreID = sales.Loja.ID
	sales.ContactID = contactID
	salesOrderJSON, err := json.Marshal(sales)
	if err != nil {
		fmt.Printf("Error marshalling sales: %v\n", err)
	}

	// Adicione um log para imprimir o JSON que está sendo enviado
	fmt.Printf("SalesOrder JSON: %s\n", string(salesOrderJSON))
	req, err := http.NewRequest("POST", baseURL+"/create_sales_order", bytes.NewBuffer(salesOrderJSON))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to create sales orders. Status: %v, Response: %s\n", resp.Status, string(body))
	}

	fmt.Printf("Sales Orders created successfully em processSales: %v\n", sales)
}

func ProcessAllItems(salesOrder types.SalesOrder) {
	for _, item := range salesOrder.Itens {
		item.SalesOrderID = salesOrder.ID
		err := ProcessItemsSalesOrder(item)
		if err != nil {
			fmt.Printf("Erro ao processar o item %d: %v\n", item.ID, err)
		}
	}
}

func ProcessItemsSalesOrder(item types.ItemsSalesOrders) error {
	// Lógica para processar e gravar o item
	err := sendItemToCreate(item)
	if err != nil {
		fmt.Printf("Erro ao criar item %d do pedido: %v\n", item.ID, err)
		return err
	}

	fmt.Printf("Item %d criado com sucesso.\n", item.ID)
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
