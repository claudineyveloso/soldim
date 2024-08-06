package types

import (
	"time"

	"github.com/google/uuid"
)

type Triage struct {
	ID                uuid.UUID `json:"id"`
	Type              string    `json:"type"`
	Grid              string    `json:"grid"`
	SkuSap            int32     `json:"sku_sap"`
	SkuWms            string    `json:"sku_wms"`
	Description       string    `json:"description"`
	CustID            int64     `json:"cust_id"`
	Seller            string    `json:"seller"`
	QuantitySupplied  int32     `json:"quantity_supplied"`
	FinalQuantity     int32     `json:"final_quantity"`
	UnitaryValue      float64   `json:"unitary_value"`
	TotalValueOffered float64   `json:"total_value_offered"`
	FinalTotalValue   float64   `json:"final_total_value"`
	Category          string    `json:"category"`
	SubCategory       string    `json:"sub_category"`
	SentToBatch       bool      `json:"sent_to_batch"`
	SentToBling       bool      `json:"sent_to_bling"`
	Defect            bool      `json:"defect"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TriageStore interface {
	ImportTriagesFromFile(filePath string) error
	CreateTriage(Triage) error
	GetTriages() ([]*Triage, error)
	GetTriageByID(id uuid.UUID) (*Triage, error)
	// UpdateTriage(Triage) error
	// DeleteTriage(id uuid.UUID) error
}
