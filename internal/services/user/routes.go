package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/claudineyveloso/soldim.git/internal/configs"
	"github.com/claudineyveloso/soldim.git/internal/services/auth"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	userStore types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	// router.HandleFunc("/register", h.handleRegister).Methods("POST")

	router.HandleFunc("/create_user", h.handleCreateUser).Methods(http.MethodPost)
	router.HandleFunc("/disabled_user", h.handleDisableUser).Methods(http.MethodPut)
	router.HandleFunc("/get_users", h.handleGetUsers).Methods(http.MethodGet)
	router.HandleFunc("/get_user/{userID}", h.handleGetUser).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.CreateLoginPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.userStore.LoginUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, expiresAt, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.LoginResponse{
		Email:     u.Email,
		IsActive:  u.IsActive,
		UserType:  u.UserType,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
// 	var user types.RegisterUserPayload
// 	if err := utils.ParseJSON(r, &user); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err := utils.Validate.Struct(user); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
// 		return
// 	}

// 	// check if user exists
// 	_, err := h.userStore.GetUserByEmail(user.Email)
// 	if err == nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
// 		return
// 	}

// 	// hash password
// 	hashedPassword, err := auth.HashPassword(user.Password)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	err = h.store.CreateUser(types.User{
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Email:     user.Email,
// 		Password:  hashedPassword,
// 	})
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusCreated, nil)
// }

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", errors))
		return
	}
	err := h.userStore.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// Define o token fixo
	// tokenFixed := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmQwZjNkMmUtYTQ3YS00NjI1LTgyMWYtZmQyODdkOThiNmZhIiwiZXhwIjoxNzI1NDk5NjIwLCJpYXQiOjE3MjU0MTMyMjB9.Z061GmVG-19Yj-Rnw9Y-_0WZ8xoN7hDibzIwgpaYcNk"

	// Extrai o token do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token de autorização não fornecido", http.StatusUnauthorized)
		return
	}

	// Extrai o token do formato "Bearer {token}"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token) // Remove possíveis espaços em branco

	if token == "" || token == authHeader {
		http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
		return
	}

	// Imprime os valores dos tokens para depuração
	// fmt.Printf("Valor do tokenFixed: %q\n", tokenFixed)
	// fmt.Printf("Valor do token recebido: %q\n", token)

	// Compara os tokens
	// if token != tokenFixed {
	// 	http.Error(w, "Token de autorização inválido", http.StatusUnauthorized)
	// 	return
	// }

	// Valida o token
	isValid, err := auth.ValidateToken(token)
	if err != nil || !isValid {
		http.Error(w, "Token de autorização inválido", http.StatusUnauthorized)
		return
	}

	// Obtém os usuários
	users, err := h.userStore.GetUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter usuários: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Define o token fixo
	// tokenFixed := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmQwZjNkMmUtYTQ3YS00NjI1LTgyMWYtZmQyODdkOThiNmZhIiwiZXhwIjoxNzI1NDk5NjIwLCJpYXQiOjE3MjU0MTMyMjB9.Z061GmVG-19Yj-Rnw9Y-_0WZ8xoN7hDibzIwgpaYcNk"

	// Extrai o token do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token de autorização não fornecido", http.StatusUnauthorized)
		return
	}

	// Extrai o token do formato "Bearer {token}"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token) // Remove possíveis espaços em branco

	if token == "" || token == authHeader {
		http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
		return
	}

	// Compara os tokens
	// if token != tokenFixed {
	// 	http.Error(w, "Token de autorização inválido", http.StatusUnauthorized)
	// 	return
	// }

	// Valida o token
	isValid, err := auth.ValidateToken(token)
	if err != nil || !isValid {
		http.Error(w, "Token de autorização inválido", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Usuário ausente"))
		return
	}
	userID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Usuário inválido"))
		return
	}

	user, err := h.userStore.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleDisableUser(w http.ResponseWriter, r *http.Request) {
	var user types.DisableUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(user); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", validationErrors))
		} else {
			utils.WriteError(w, http.StatusBadRequest, err)
		}
		return
	}
	err := h.userStore.DisableUser(r.Context(), user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
