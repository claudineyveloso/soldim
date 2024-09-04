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

type UserStore interface {
	GetUserByID(userID uuid.UUID) (*types.User, error)
}

// WithJWTAuth middleware que valida o token JWT e adiciona o usuário ao contexto
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtém o token do request
		tokenString := utils.GetTokenFromRequest(r)
		token, err := ValidateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		// Extrai as claims do token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("invalid token claims")
			permissionDenied(w)
			return
		}

		// Obtém o userID da claim
		str, ok := claims["userID"].(string)
		if !ok {
			log.Println("userID claim is missing or invalid")
			permissionDenied(w)
			return
		}

		// Converte o userID para UUID
		userID, err := uuid.Parse(str)
		if err != nil {
			log.Printf("failed to parse userID: %v", err)
			permissionDenied(w)
			return
		}

		// Obtém o usuário do store
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Adiciona o usuário ao contexto
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// Chama o handler original
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID uuid.UUID, sessionVersion string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24) // Token expira em 24 horas
	// Define as claims do token
	claims := &types.Claims{
		UserID:         userID,
		SessionVersion: sessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // Token expira em 24 horas
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

func ValidateToken(tokenString string, userStore UserStore) (bool, error) {
	// Defina a estrutura das claims
	claims := &types.Claims{}

	// Parse o token com claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifique o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(configs.Envs.JWTSecret), nil
	})
	if err != nil {
		fmt.Printf("Erro ao decodificar o token: %v\n", err)
		return false, err
	}

	// Verifique se o token é válido
	if !token.Valid {
		fmt.Println("Token inválido")
		return false, errors.New("token inválido")
	}

	// Verifique se o token expirou
	if time.Now().After(claims.ExpiresAt.Time) {
		fmt.Println("Token expirado")
		return false, errors.New("token expirado")
	}

	// Verifique a versão da sessão
	user, err := userStore.GetUserByID(claims.UserID)
	if err != nil {
		fmt.Printf("Erro ao obter o usuário: %v\n", err)
		return false, fmt.Errorf("erro ao obter o usuário: %v", err)
	}

	// Verifique a versão da sessão
	if claims.SessionVersion != user.SessionVersion {
		fmt.Printf("Sessão não corresponde: token tem %s, banco tem %s\n", claims.SessionVersion, user.SessionVersion)
		return false, errors.New("token não é mais válido")
	}

	return true, nil
}

// func ValidateToken(tokenString string, userStore UserStore) (bool, error) {
// 	claims := &types.Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
// 		}
// 		return []byte(configs.Envs.JWTSecret), nil
// 	})
// 	if err != nil {
// 		fmt.Printf("Erro ao decodificar o token: %v\n", err)
// 		return false, err
// 	}

// 	if !token.Valid {
// 		fmt.Println("Token inválido")
// 		return false, errors.New("token inválido")
// 	}

// 	if time.Now().After(claims.ExpiresAt.Time) {
// 		fmt.Println("Token expirado")
// 		return false, errors.New("token expirado")
// 	}

// 	// Verifique a versão da sessão
// 	user, err := userStore.GetUserByID(claims.UserID)
// 	if err != nil {
// 		return false, fmt.Errorf("erro ao obter o usuário: %v", err)
// 	}
// 	// Verifique a versão da sessão
// 	if claims.SessionVersion != user.SessionVersion {
// 		return false, errors.New("token não é mais válido")
// 	}

// 	return true, nil
// }

func ValidateToken222(tokenString string) (bool, error) {
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

	fmt.Printf("Token válido: %v\n", token.Valid)
	if !token.Valid {
		fmt.Println("Token inválido")
		return false, errors.New("token inválido")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		fmt.Println("Token expirado")
		return false, errors.New("token expirado")
	}

	fmt.Println("Token válido e não expirado")
	return true, nil
}

func ValidateToken123(tokenString string) (bool, error) {
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

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Define a estrutura das claims
	claims := &types.Claims{}

	// Analisa o token com as claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifique o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		// Retorna a chave secreta usada para assinatura
		return []byte(configs.Envs.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar o token: %v", err)
	}

	// Verifique se o token é válido
	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Verifique se o token expirou
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expirado")
	}

	return token, nil
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
