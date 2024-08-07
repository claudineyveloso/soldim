package types

import "time"

type Contact struct {
	ID              int64     `json:"id"`
	Nome            string    `json:"nome"`
	Codigo          string    `json:"codigo"`
	Situacao        string    `json:"situacao"`
	Numerodocumento string    `json:"numero_documento"`
	Telefone        string    `json:"telefone"`
	Celular         string    `json:"celular"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ContactResponse struct {
	Data []Contact `json:"data"`
}

type ContactStore interface {
	CreateContact(Contact) error
	GetContacts() ([]*Contact, error)
}
