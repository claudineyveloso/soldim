package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims define a estrutura do payload do JWT
type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	SessionVersion string    `json:"session_version"`
	jwt.RegisteredClaims
}
