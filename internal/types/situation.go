package types

import "time"

type Situation struct {
	ID        int64     `json:"id"`
	Descricao string    `json:"descricao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SituationStore interface {
	CreateSituation(Situation) error
	GetSituations() ([]*Situation, error)
	GetSituationByID(id int64) (*Situation, error)
	GetSituationByDescription(descricao string) (*Situation, error)
}
