package contact

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateContact(contact types.Contact) error {
	queries := db.New(s.db)
	ctx := context.Background()

	if contact.ID == 0 {
		return fmt.Errorf("invalid contact ID: %d", contact.ID)
	}

	now := time.Now()
	contact.CreatedAt = now
	contact.UpdatedAt = now

	createContactParams := db.CreateContactParams{
		ID:              contact.ID,
		Nome:            contact.Nome,
		Codigo:          contact.Codigo,
		Situacao:        contact.Situacao,
		Numerodocumento: contact.Numerodocumento,
		Telefone:        contact.Telefone,
		Celular:         contact.Celular,
		CreatedAt:       contact.CreatedAt,
		UpdatedAt:       contact.UpdatedAt,
	}

	if err := queries.CreateContact(ctx, createContactParams); err != nil {
		fmt.Println("Erro ao criar um Contato:", err)
		return err
	}
	fmt.Printf("Contato criado com ID: %d\n", contact.ID)
	return nil
}

func (s *Store) GetContacts() ([]*types.Contact, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbContacts, err := queries.GetContacts(ctx)
	if err != nil {
		return nil, err
	}

	var contacts []*types.Contact
	for _, dbContact := range dbContacts {
		contact := convertDBContactToContact(dbContact)
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (s *Store) GetContactByID(contactID int64) (*types.Contact, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	// Tenta buscar o contato no banco de dados
	dbContact, err := queries.GetContact(ctx, contactID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se o contato não for encontrado, retorna nil sem erro
			return &types.Contact{}, nil
		}
		// Retorna o erro se houver algum problema na consulta ao banco de dados
		return nil, err
	}

	// Converte o contato do formato do banco de dados para o formato da aplicação
	contact := convertDBContactToContact(dbContact)

	// Retorna o contato encontrado
	return contact, nil
}

func convertDBContactToContact(dbContact db.Contact) *types.Contact {
	contact := &types.Contact{
		ID:              dbContact.ID,
		Nome:            dbContact.Nome,
		Codigo:          dbContact.Codigo,
		Situacao:        dbContact.Situacao,
		Numerodocumento: dbContact.Numerodocumento,
		Telefone:        dbContact.Telefone,
		Celular:         dbContact.Celular,
		CreatedAt:       dbContact.CreatedAt,
		UpdatedAt:       dbContact.UpdatedAt,
	}
	return contact
}
