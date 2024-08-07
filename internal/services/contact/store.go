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
