package token

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	tokenStore types.TokenStore
}

func NewHandler(tokenStore types.TokenStore) *Handler {
	return &Handler{tokenStore: tokenStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_token", h.handleCreateToken).Methods(http.MethodPost)
	router.HandleFunc("/get_token", h.handleGetToken).Methods(http.MethodGet)
	router.HandleFunc("/update_token", h.handleUpdateToken).Methods(http.MethodPut)
	router.HandleFunc("/delete_token/{tokenID}", h.handleDeleteToken).Methods(http.MethodDelete)
}

func (h *Handler) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	var token types.TokenPayload
	if err := utils.ParseJSON(r, &token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(token); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.tokenStore.CreateToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)
}

func (h *Handler) handleUpdateToken(w http.ResponseWriter, r *http.Request) {
	var token types.TokenPayload
	if err := utils.ParseJSON(r, &token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(token); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.tokenStore.UpdateToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    token,
		"message": "Registro alterado com sucesso",
		"status":  http.StatusOK,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func (h *Handler) handleGetToken(w http.ResponseWriter, r *http.Request) {
	// bucketID := auth.GetUserIDFromContext(r.Context())
	// fmt.Println("Valor de userIDffsadfsda", bucketID)
	token, err := h.tokenStore.GetToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Token: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, token)
}

func (h *Handler) handleDeleteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["tokenID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Token ausente!"))
		return
	}
	parsedTokenID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Token inválido!"))
		return
	}

	err = h.tokenStore.DeleteToken(parsedTokenID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "Registro apagado com sucesso",
		"status":  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}
