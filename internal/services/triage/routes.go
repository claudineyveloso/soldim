package triage

import (
	"encoding/json"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/gorilla/mux"
)

type Handler struct {
	triageStore types.TriageStore
}

func NewHandler(triageStore types.TriageStore) *Handler {
	return &Handler{triageStore: triageStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_triage", h.handleCreateTriage).Methods(http.MethodPost)
	router.HandleFunc("/get_triages", h.handleGetTriages).Methods(http.MethodGet)
	router.HandleFunc("/import_triages", h.handleImportTriage).Methods(http.MethodGet)
}

func (h *Handler) handleCreateTriage(w http.ResponseWriter, r *http.Request) {
	var t types.Triage
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	if err := h.triageStore.CreateTriage(t); err != nil {
		http.Error(w, "erro ao criar triagem", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleGetTriages(w http.ResponseWriter, r *http.Request) {
	triages, err := h.triageStore.GetTriages()
	if err != nil {
		http.Error(w, "erro ao buscar triagens", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(triages); err != nil {
		http.Error(w, "erro ao codificar JSON", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleImportTriage(w http.ResponseWriter, r *http.Request) {
	filePath := "internal/files/LOTE_188.xlsx"
	if err := h.triageStore.ImportTriagesFromFile(filePath); err != nil {
		http.Error(w, "erro ao importar triagens: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}