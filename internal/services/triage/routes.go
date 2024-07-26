package triage

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
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
	description := r.URL.Query().Get("description")
	sku_wms := r.URL.Query().Get("sku_wms")
	sku_sapStr := r.URL.Query().Get("sku_sap")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	var sku_sap int32
	if sku_sapStr != "" {
		sku_sapInt, err := strconv.Atoi(sku_sapStr)
		if err != nil {
			http.Error(w, "invalid sku_sap value", http.StatusBadRequest)
			return
		}
		sku_sap = int32(sku_sapInt)
	}

	triages, totalCount, err := h.triageStore.GetTriages(description, sku_wms, sku_sap, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, "erro ao buscar triagens", http.StatusInternalServerError)
		return
	}

	response := struct {
		Triages    []*types.Triage `json:"triage"`
		TotalCount int64           `json:"total_count"`
	}{
		Triages:    triages,
		TotalCount: totalCount,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleImportTriage(w http.ResponseWriter, r *http.Request) {
	filePath := "internal/files/LOTE_188.xlsx"
	if err := h.triageStore.ImportTriagesFromFile(filePath); err != nil {
		http.Error(w, "erro ao importar triagens: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
