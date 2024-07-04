package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// Helper function to handle HTTP errors
func HandleHTTPError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

// CreateNullString returns a sql.NullString with the given value, if not empty
// func CreateNullString(value string) sql.NullString {
// 	if value != "" {
// 		return sql.NullString{String: value, Valid: true}
// 	}
// 	return sql.NullString{Valid: false}
// }

func CreateNullString(value interface{}) sql.NullString {
	switch v := value.(type) {
	case string:
		if v != "" {
			return sql.NullString{String: v, Valid: true}
		}
	case sql.NullString:
		return v
	}
	return sql.NullString{Valid: false}
}

func CreateNullDate(value string) sql.NullTime {
	// Tente analisar a string de data
	parsedTime, err := time.Parse("2006-01-02", value)
	if err != nil {
		// Se houver um erro ao analisar a data, retorne uma NullTime inválida
		return sql.NullTime{Valid: false}
	}
	// Se a data for válida, retorne uma NullTime com o tempo definido
	return sql.NullTime{Time: parsedTime, Valid: true}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func Uint32ToUUIDBytes(id uint32) []byte {
	bytes := make([]byte, 16)
	bytes[12] = byte(id >> 24)
	bytes[13] = byte(id >> 16)
	bytes[14] = byte(id >> 8)
	bytes[15] = byte(id)
	return bytes
}

func GetValidString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	if nullStr, ok := value.(sql.NullString); ok && nullStr.Valid {
		return nullStr.String
	}
	return ""
}

// CalculateRelevanceScore calcula a pontuação de relevância de uma descrição de produto com base na consulta de pesquisa.
func CalculateRelevanceScore(query, description string) int {
	queryWords := strings.Fields(strings.ToLower(query))
	descriptionWords := strings.Fields(strings.ToLower(description))

	score := 0
	for _, queryWord := range queryWords {
		for _, descWord := range descriptionWords {
			if queryWord == descWord {
				score++
			}
		}
	}
	return score
}
