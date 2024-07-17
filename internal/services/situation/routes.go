package situation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	situationStore types.SituationStore
}

func NewHandler(situationStore types.SituationStore) *Handler {
	return &Handler{situationStore: situationStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_situation", h.handleCreateSituation).Methods(http.MethodPost)
	router.HandleFunc("/get_situations", h.handleGetSituations).Methods(http.MethodGet)
	router.HandleFunc("/get_situation/{situationID}", h.handleGetSituation).Methods(http.MethodGet)
	router.HandleFunc("/get_situation_description/{description}", h.handleGetSituationByDescription).Methods(http.MethodGet)
}

func (h *Handler) handleCreateSituation(w http.ResponseWriter, r *http.Request) {
	var situation types.Situation
	if err := utils.ParseJSON(r, &situation); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(situation); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.situationStore.CreateSituation(situation)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    situation,
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

func (h *Handler) handleGetSituations(w http.ResponseWriter, r *http.Request) {
	draft, err := h.situationStore.GetSituations()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Rascunho: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, draft)
}

func (h *Handler) handleGetSituation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	situationIDStr, ok := vars["situationID"]
	if !ok {
		http.Error(w, "ID da Situação ausente!", http.StatusBadRequest)
		return
	}

	situationID, err := strconv.ParseInt(situationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID da Situação inválido!", http.StatusBadRequest)
		return
	}

	situation, err := h.situationStore.GetSituationByID(situationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter a Situação: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, situation)
}

func (h *Handler) handleGetSituationByDescription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	description, ok := vars["description"]
	if !ok {
		http.Error(w, "Descrição da Situação ausente!", http.StatusBadRequest)
		return
	}

	situation, err := h.situationStore.GetSituationByDescription(description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter a Situação: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, situation)
}
