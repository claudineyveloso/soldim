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

func (s *Store) CreateToken(token types.TokenPayload) error {
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
	}

	if err := queries.CreateToken(ctx, createTokenParams); err != nil {
		// http.Error(_, "Erro ao criar usuário", http.StatusInternalServerError)
		fmt.Println("Erro ao criar um Bucket:", err)
		return err
	}
	return nil
}

func (s *Store) UpdateToken(token types.TokenPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	updateTokenParams := db.UpdateTokenParams{
		ID:           token.ID,
		AccessToken:  token.AccessToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
		Scope:        token.Scope,
		RefreshToken: token.RefreshToken,
	}

	if err := queries.UpdateToken(ctx, updateTokenParams); err != nil {
		fmt.Println("Erro ao atualizar um Token:", err)
		return err
	}
	return nil
}

func (s *Store) DeleteToken(tokenID uuid.UUID) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteToken(ctx, tokenID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetToken() ([]*types.Token, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbTokens, err := queries.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	var tokens []*types.Token
	for _, dbToken := range dbTokens {
		token := convertDBTokenToToken(dbToken)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func convertDBTokenToToken(dbToken db.Token) *types.Token {
	token := &types.Token{
		ID:           dbToken.ID,
		AccessToken:  dbToken.AccessToken,
		ExpiresIn:    dbToken.ExpiresIn,
		TokenType:    dbToken.TokenType,
		Scope:        dbToken.Scope,
		RefreshToken: dbToken.RefreshToken,
	}
	return token
}
