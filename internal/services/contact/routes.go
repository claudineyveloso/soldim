package contact

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
	contactStore types.ContactStore
}

func NewHandler(contactStore types.ContactStore) *Handler {
	return &Handler{contactStore: contactStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_contact", h.handleCreateContacts).Methods(http.MethodPost)
	router.HandleFunc("/get_contacts", h.handleGetContacts).Methods(http.MethodGet)
}

func (h *Handler) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.contactStore.GetContacts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os contatos : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, contacts)
}

func (h *Handler) handleCreateContacts(w http.ResponseWriter, r *http.Request) {
	var contact types.Contact
	if err := utils.ParseJSON(r, &contact); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(contact); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.contactStore.CreateContact(contact)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    contact,
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
