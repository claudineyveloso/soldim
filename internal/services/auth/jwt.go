package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "userID"

// Defina sua chave secreta (deve ser mantida segura e não compartilhada)
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims define a estrutura do payload do JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		// userID, err := strconv.Atoi(str)
		// if err != nil {
		// 	log.Printf("failed to convert userID to int: %v", err)
		// 	permissionDenied(w)
		// 	return
		// }

		userID, err := uuid.Parse(str)
		if err != nil {
			log.Printf("failed to parse userID: %v", err)
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

// Função para criar o JWT
func CreateJWT(secret []byte, userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24) // Token expira em 24 horas
	// Define as claims do token
	claims := &types.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expira em 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Cria o token usando o método de assinatura HS256 e as claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func ValidateToken(tokenString string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(configs.Envs.JWTSecret), nil
	})
	if err != nil {
		fmt.Printf("Erro ao decodificar o token: %v\n", err)
		return false, err
	}

	if !token.Valid {
		fmt.Println("Token inválido")
		return false, errors.New("token inválido")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		fmt.Println("Token expirado")
		return false, errors.New("token expirado")
	}

	return true, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
