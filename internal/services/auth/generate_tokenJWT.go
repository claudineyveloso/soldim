package auth

import (
	"time"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("minha_chave_secreta") // Use uma chave mais segura

// Função para gerar um token JWT
func GenerateJWT(payload types.TokenPayload) (string, error) {
	claims := &jwt.MapClaims{
		"sub": payload.ID,                            // Supondo que TokenPayload tem um campo UserID
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token válido por 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
