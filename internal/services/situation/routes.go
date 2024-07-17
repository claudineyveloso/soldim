package situation

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
	situatiomnStore types.SituationStore
}

func NewHandler(situationStore types.SituationStore) *Handler {
	return &Handler{situatiomnStore: situationStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_situation", h.handleCreateSituation).Methods(http.MethodPost)
	router.HandleFunc("/get_situations", h.handleGetSituations).Methods(http.MethodGet)
	router.HandleFunc("/get_situation/{situationID}", h.handleGetSituation).Methods(http.MethodGet)
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
	err := h.situatiomnStore.CreateSituation(situation)
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
