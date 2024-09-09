package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func ProcessContact(sale types.SalesOrder) int64 {
	contact, err := ExistContact(sale.Contato.ID)
	if err != nil {
		fmt.Printf("Erro verificando se o contato existe ID %d: %v\n", sale.Contato.ID, err)
		return 0
	}

	if contact == nil || contact.ID == 0 {
		// Conversão do tipo Contato para Contact
		newContact := types.Contato{
			ID:              sale.Contato.ID,
			Nome:            sale.Contato.Nome,
			TipoPessoa:      sale.Contato.TipoPessoa,
			NumeroDocumento: sale.Contato.NumeroDocumento,
		}
		contact, err = CreateContact(newContact)
		if err != nil {
			fmt.Printf("Erro ao criar contato ID %d: %v\n", sale.Contato.ID, err)
			return 0
		}
		sale.Contato.ID = contact.ID
	} else {
		sale.Contato.ID = contact.ID
		fmt.Printf("Contato ID %d já existe\n", sale.Contato.ID)
	}
	sale.Contato.ID = contact.ID
	sale.ContactID = contact.ID
	return contact.ID
}

func CreateContact(contato types.Contato) (*types.Contato, error) {
	url := baseURL + "/create_contact"
	contactData, err := json.Marshal(contato)
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

	var createdContact types.Contato
	err = json.NewDecoder(resp.Body).Decode(&createdContact)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	// Verifica se o contato realmente existe no banco de dados após a criação
	existingContact, err := ExistContact(createdContact.ID)
	if err != nil {
		return nil, fmt.Errorf("contact not found after creation: %v", err)
	}

	return existingContact, nil
}
