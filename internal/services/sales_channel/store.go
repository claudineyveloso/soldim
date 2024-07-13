package saleschannel

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

func (s *Store) CreateSalesChannel(saleschannel types.SalesChannel) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	saleschannel.CreatedAt = now
	saleschannel.UpdatedAt = now

	createSalesChannelParams := db.CreateSalesChannelParams{
		ID:        saleschannel.ID,
		Descricao: saleschannel.Descricao,
		Tipo:      saleschannel.Tipo,
		Situacao:  saleschannel.Situacao,
		CreatedAt: saleschannel.CreatedAt,
		UpdatedAt: saleschannel.UpdatedAt,
	}

	if err := queries.CreateSalesChannel(ctx, createSalesChannelParams); err != nil {
		fmt.Println("Erro ao criar um canal de vendas:", err)
		return err
	}
	return nil
}

func (s *Store) GetSalesChannel() ([]*types.SalesChannel, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbSalesChannels, err := queries.GetSalesChannel(ctx)
	if err != nil {
		return nil, err
	}

	var saleschannels []*types.SalesChannel
	for _, dbSalesChannel := range dbSalesChannels {
		saleschannel := convertDBSalesChannelToSalesChannel(dbSalesChannel)
		saleschannels = append(saleschannels, saleschannel)
	}
	return saleschannels, nil
}

func convertDBSalesChannelToSalesChannel(dbSalesChannel db.SalesChannel) *types.SalesChannel {
	saleschannel := &types.SalesChannel{
		ID:        dbSalesChannel.ID,
		Descricao: dbSalesChannel.Descricao,
		Tipo:      dbSalesChannel.Tipo,
		Situacao:  dbSalesChannel.Situacao,
		CreatedAt: dbSalesChannel.CreatedAt,
		UpdatedAt: dbSalesChannel.UpdatedAt,
	}
	return saleschannel
}
