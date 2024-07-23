package triage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ImportTriagesFromFile(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("erro ao obter linhas da planilha: %v", err)
	}

	var triages []*types.Triage
	for _, row := range rows[1:] { // Pular o cabeçalho
		id, err := uuid.Parse(row[0])
		if err != nil {
			fmt.Printf("Erro ao parsear UUID: %v\n", err)
			continue
		}
		triage := &types.Triage{
			ID:                id,
			Type:              row[1],
			Grid:              row[2],
			SkuSap:            parseInt32(row[3]),
			SkuWms:            row[4],
			Description:       row[5],
			CustID:            parseInt64(row[6]),
			Seller:            row[7],
			QuantitySupplied:  parseInt32(row[8]),
			FinalQuantity:     parseInt32(row[9]),
			UnitaryValue:      parseFloat(row[10]),
			TotalValueOffered: parseFloat(row[11]),
			FinalTotalValue:   parseFloat(row[12]),
			Category:          row[13],
			SubCategory:       row[14],
			SentToBatch:       parseBool(row[15]),
			SentToBling:       parseBool(row[16]),
			Defect:            parseBool(row[17]),
		}
		triages = append(triages, triage)
	}

	return s.ImportTriages(triages)
}

func (s *Store) ImportTriages(triages []*types.Triage) error {
	queries := db.New(s.db)
	ctx := context.Background()
	for _, triage := range triages {
		now := time.Now()
		triage.CreatedAt = now
		triage.UpdatedAt = now
		createTriageParams := db.CreateTriageParams{
			ID:                triage.ID,
			Type:              triage.Type,
			Grid:              triage.Grid,
			SkuSap:            triage.SkuSap,
			SkuWms:            triage.SkuWms,
			Description:       triage.Description,
			CustID:            triage.CustID,
			Seller:            triage.Seller,
			QuantitySupplied:  triage.QuantitySupplied,
			FinalQuantity:     triage.FinalQuantity,
			UnitaryValue:      triage.UnitaryValue,
			TotalValueOffered: triage.TotalValueOffered,
			FinalTotalValue:   triage.FinalTotalValue,
			Category:          triage.Category,
			SubCategory:       triage.SubCategory,
			SentToBatch:       triage.SentToBatch,
			SentToBling:       triage.SentToBling,
			Defect:            triage.Defect,
			CreatedAt:         triage.CreatedAt,
			UpdatedAt:         triage.UpdatedAt,
		}
		if err := queries.CreateTriage(ctx, createTriageParams); err != nil {
			fmt.Println("Erro ao criar um triagem:", err)
			return err
		}
	}
	return nil
}

func (s *Store) GetTriages() ([]*types.Triage, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbTriages, err := queries.GetTriages(ctx)
	if err != nil {
		return nil, err
	}

	var triages []*types.Triage
	for _, dbTriage := range dbTriages {
		triage := convertDBTriageToTriage(dbTriage)
		triages = append(triages, triage)
	}
	return triages, nil
}

func parseInt32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func (s *Store) CreateTriage(triage types.Triage) error {
	queries := db.New(s.db)
	ctx := context.Background()
	now := time.Now()
	triage.CreatedAt = now
	triage.UpdatedAt = now
	createTriageParams := db.CreateTriageParams{
		ID:                triage.ID,
		Type:              triage.Type,
		Grid:              triage.Grid,
		SkuSap:            triage.SkuSap,
		SkuWms:            triage.SkuWms,
		Description:       triage.Description,
		CustID:            triage.CustID,
		Seller:            triage.Seller,
		QuantitySupplied:  triage.QuantitySupplied,
		FinalQuantity:     triage.FinalQuantity,
		UnitaryValue:      triage.UnitaryValue,
		TotalValueOffered: triage.TotalValueOffered,
		FinalTotalValue:   triage.FinalTotalValue,
		Category:          triage.Category,
		SubCategory:       triage.SubCategory,
		SentToBatch:       triage.SentToBatch,
		SentToBling:       triage.SentToBling,
		Defect:            triage.Defect,
		CreatedAt:         triage.CreatedAt,
		UpdatedAt:         triage.UpdatedAt,
	}
	if err := queries.CreateTriage(ctx, createTriageParams); err != nil {
		fmt.Println("Erro ao criar um triagem:", err)
		return err
	}
	return nil
}

func convertDBTriageToTriage(dbTriage db.Triage) *types.Triage {
	triage := &types.Triage{
		ID:                dbTriage.ID,
		Type:              dbTriage.Type,
		Grid:              dbTriage.Grid,
		SkuSap:            dbTriage.SkuSap,
		SkuWms:            dbTriage.SkuWms,
		Description:       dbTriage.Description,
		CustID:            dbTriage.CustID,
		Seller:            dbTriage.Seller,
		QuantitySupplied:  dbTriage.QuantitySupplied,
		FinalQuantity:     dbTriage.FinalQuantity,
		UnitaryValue:      dbTriage.UnitaryValue,
		TotalValueOffered: dbTriage.TotalValueOffered,
		FinalTotalValue:   dbTriage.FinalTotalValue,
		Category:          dbTriage.Category,
		SubCategory:       dbTriage.SubCategory,
		SentToBatch:       dbTriage.SentToBatch,
		SentToBling:       dbTriage.SentToBling,
		Defect:            dbTriage.Defect,
		CreatedAt:         dbTriage.CreatedAt,
		UpdatedAt:         dbTriage.UpdatedAt,
	}
	return triage
}
