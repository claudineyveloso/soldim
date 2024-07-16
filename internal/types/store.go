package types

import "time"

type Store struct {
	ID        int64     `json:"id"`
	Descricao string    `json:"descricao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StoreStore interface {
	CreateStore(Store) error
	UpdateStore(Store) error
}
