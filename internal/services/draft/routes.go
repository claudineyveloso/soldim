package draft

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
	draftStore types.DraftStore
}

func NewHandler(draftStore types.DraftStore) *Handler {
	return &Handler{draftStore: draftStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_draft", h.handleCreateDraft).Methods(http.MethodPost)
	router.HandleFunc("/get_drafts", h.handleGetDrafts).Methods(http.MethodGet)
	router.HandleFunc("/get_draft/{draftID}", h.handleGetDraft).Methods(http.MethodGet)
	router.HandleFunc("/get_drafts_by_search_id/{searchID}", h.handleGetDraftBySearchID).Methods(http.MethodGet)
	router.HandleFunc("/update_draft", h.handleUpdateDraft).Methods(http.MethodPut)
	router.HandleFunc("/delete_draft/{draftID}", h.handleDeleteDraft).Methods(http.MethodDelete)
	router.HandleFunc("/delete_drafts_by_search_id/{searchID}", h.handleDeleteDraftBySearchID).Methods(http.MethodDelete)
}

func (h *Handler) handleGetDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.draftStore.GetDrafts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os Rascunhos : %v", err), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		Drafts []*types.Draft `json:"drafts"`
	}{Drafts: drafts})
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao serializar os produtos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) handleCreateDraft(w http.ResponseWriter, r *http.Request) {
	var draft types.DraftPayload
	if err := utils.ParseJSON(r, &draft); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(draft); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.draftStore.CreateDraft(draft)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    draft,
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

func (h *Handler) handleDeleteDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["draftID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Draft ausente!"))
		return
	}
	parsedDraftsID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Draft inválido!"))
		return
	}

	err = h.draftStore.DeleteDraft(parsedDraftsID)
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

func (h *Handler) handleDeleteDraftBySearchID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Search ausente!"))
		return
	}
	parsedSearchesID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Search inválido!"))
		return
	}

	err = h.draftStore.DeleteDraftBySearchID(parsedSearchesID)
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

func (h *Handler) handleGetDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["draftID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Rascunho ausente!"))
		return
	}
	parsedDraftsID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Rascunho inválido!"))
		return
	}

	draft, err := h.draftStore.GetDraftByID(parsedDraftsID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, draft)
}

func (h *Handler) handleGetDraftBySearchID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Search ausente!"))
		return
	}
	parsedSearchID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Search inválido!"))
		return
	}

	draft, err := h.draftStore.GetDraftBySearchID(parsedSearchID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, draft)
}

func (h *Handler) handleUpdateDraft(w http.ResponseWriter, r *http.Request) {
	var draft types.DraftPayload
	if err := utils.ParseJSON(r, &draft); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(draft); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.draftStore.UpdateDraft(draft)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    draft,
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
