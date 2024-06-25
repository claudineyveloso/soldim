package token

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Createtoken(token types.TokenPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	token.ID = uuid.New()
	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now

	createTokenParams := db.CreateTokenParams{
		ID:           token.ID,
		AccessToken:  token.AccessToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
		Scope:        token.Scope,
		RefreshToken: token.RefreshToken,
		CreatedAt:    token.CreatedAt,
		UpdatedAt:    token.UpdatedAt,
	}

	if err := queries.CreateToken(ctx, createTokenParams); err != nil {
		// http.Error(_, "Erro ao criar usuário", http.StatusInternalServerError)
		fmt.Println("Erro ao criar um Bucket:", err)
		return err
	}
	return nil
}
