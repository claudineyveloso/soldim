package deposit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	depositStore types.DepositStore
}

func NewHandler(depositStore types.DepositStore) *Handler {
	return &Handler{depositStore: depositStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_deposit", h.handleCreateDeposits).Methods(http.MethodPost)
	router.HandleFunc("/get_deposits", h.handleGetDeposits).Methods(http.MethodGet)
}

func (h *Handler) handleGetDeposits(w http.ResponseWriter, r *http.Request) {
	deposits, err := h.depositStore.GetDeposits()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os depositos : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, deposits)
}

func (h *Handler) handleCreateDeposits(w http.ResponseWriter, r *http.Request) {
	var deposit types.Deposit
	if err := utils.ParseJSON(r, &deposit); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(deposit); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.depositStore.CreateDeposit(deposit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    deposit,
		"message": "Registro criado com sucesso",
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
