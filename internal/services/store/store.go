package store

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

func (s *Store) CreateStore(store types.Store) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	store.CreatedAt = now
	store.UpdatedAt = now

	createStoreParams := db.CreateStoreParams{
		ID:        store.ID,
		Descricao: store.Descricao,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}

	if err := queries.CreateStore(ctx, createStoreParams); err != nil {
		fmt.Println("Erro ao criar uma loja:", err)
		return err
	}
	return nil
}

func (s *Store) GetStores() ([]*types.Store, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbStores, err := queries.GetStores(ctx)
	if err != nil {
		return nil, err
	}

	var stores []*types.Store
	for _, dbStore := range dbStores {
		store := convertDBStoreToStore(dbStore)
		stores = append(stores, store)
	}
	return stores, nil
}

func (s *Store) GetStoreByID(storeID int64) (*types.Store, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbStore, err := queries.GetStore(ctx, storeID)
	if err != nil {
		return nil, err
	}
	store := convertDBStoreToStore(dbStore)

	return store, nil
}

func (s *Store) GetStoreByDescription(description string) (*types.Store, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbStore, err := queries.GetStoreByDescription(ctx, description)
	if err != nil {
		return nil, err
	}
	product := convertDBStoreToStore(dbStore)

	return product, nil
}

func convertDBStoreToStore(dbStore db.Store) *types.Store {
	store := &types.Store{
		ID:        dbStore.ID,
		Descricao: dbStore.Descricao,
		CreatedAt: dbStore.CreatedAt,
		UpdatedAt: dbStore.UpdatedAt,
	}
	return store
}
