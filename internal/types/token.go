package types

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int32     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TokenPayload struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int32     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TokenStore interface {
	CreateToken(TokenPayload) error
	GetToken() ([]*Token, error)
	UpdateToken(TokenPayload) error
	DeleteToken(id uuid.UUID) error
}
