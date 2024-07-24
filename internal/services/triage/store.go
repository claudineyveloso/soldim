package triage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
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
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo não encontrado: %s", filePath)
	}
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer f.Close()

	// Listar todas as planilhas
	sheets := f.GetSheetList()
	fmt.Printf("Planilhas disponíveis: %v\n", sheets)

	// Verificar se "Sheet1" existe
	sheetName := "BRMG01"
	if !contains(sheets, sheetName) {
		return fmt.Errorf("a planilha %s não existe. Planilhas disponíveis: %v", sheetName, sheets)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("erro ao obter linhas da planilha: %v", err)
	}
	var triages []*types.Triage
	for i, row := range rows { // Loop through all rows
		if i < 3 {
			// Pular a primeira linha (cabeçalho)
			continue
		}
		if len(row) < 14 { // Verificar se a linha tem pelo menos 14 colunas
			log.Printf("Linha %d com dados insuficientes: %v\n", i+1, row)
			continue
		}

		// Adicionar log para verificar o número de colunas
		log.Printf("Processando linha %d com %d colunas\n", i+1, len(row))

		triage := &types.Triage{
			Type:              row[0],
			Grid:              row[1],
			SkuSap:            parseInt32(row[2]),
			SkuWms:            row[3],
			Description:       row[4],
			CustID:            parseInt64(row[5]),
			Seller:            row[6],
			QuantitySupplied:  parseInt32(row[7]),
			FinalQuantity:     parseInt32(row[8]),
			UnitaryValue:      parseFloat(row[9]),
			TotalValueOffered: parseFloat(row[10]),
			FinalTotalValue:   parseFloat(row[11]),
			Category:          row[12],
			SubCategory:       row[13],
		}
		triages = append(triages, triage)
	}

	return s.ImportTriages(triages)
}

// Função para verificar se uma planilha existe na lista
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (s *Store) ImportTriages(triages []*types.Triage) error {
	queries := db.New(s.db)
	ctx := context.Background()
	for _, triage := range triages {
		id := uuid.New()

		now := time.Now()
		triage.CreatedAt = now
		triage.UpdatedAt = now
		triage.SentToBatch = false
		triage.SentToBling = false
		triage.Defect = false
		createTriageParams := db.CreateTriageParams{
			ID:                id,
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
