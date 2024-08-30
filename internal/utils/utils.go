package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type TokenResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

var (
	Validate = validator.New()
	baseURL  = GetBaseURL()
)

// Retorna a URL base dependendo do ambiente
func GetBaseURL() string {
	env := os.Getenv("ENVIRONMENT")

	if env == "prod" {
		return "https://soldim-api-344942c665db.herokuapp.com"
	}
	return "http://localhost:8080"
}

// Helper function to handle HTTP errors
func HandleHTTPError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

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

	return json.NewDecoder(r.Body).Decode(&v)
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

func LogError(logFile *os.File, ID int64, err error) {
	logEntry := fmt.Sprintf("ID: %d, Error: %v\n", ID, err)
	if _, writeErr := logFile.WriteString(logEntry); writeErr != nil {
		fmt.Printf("Erro ao escrever no arquivo de log: %v\n", writeErr)
	}
}

func LogErrorToFile(logMessage string) {
	// Define the log file path
	logFilePath := "errors.log"

	// Open the log file in append mode, create it if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Create a timestamp for the log entry
	timestamp := time.Now().Format(time.RFC3339)

	// Write the log message to the file
	_, err = fmt.Fprintf(logFile, "[%s] %s", timestamp, logMessage)
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

func FetchAccessToken() (string, error) {
	// URL da API para obter o token de acesso
	url := baseURL + "/get_token"

	// Envia a requisição para obter o token
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("resposta da API não OK: %s", resp.Status)
	}

	// Decodifica a resposta JSON como um array de TokenResponse
	var responseArray []TokenResponse
	if err := json.Unmarshal(body, &responseArray); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	// Verifica se o array tem pelo menos um item
	if len(responseArray) < 1 {
		return "", fmt.Errorf("resposta inesperada: array vazio")
	}

	// Retorna o access_token do primeiro item no array
	return responseArray[0].AccessToken, nil
}

func ParseFloat(value string) float64 {
	// Verificar se existe "R$" no valor e removê-lo se estiver presente
	if strings.Contains(value, "R$") {
		value = strings.Replace(value, "R$", "", -1)
	}

	// Remover espaços em branco e substituir "," por "."
	value = strings.TrimSpace(strings.Replace(value, ",", ".", -1))

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Erro ao converter %s para float64: %v\n", value, err)
		return 0
	}
	return floatValue
}

func ParseCurrency(value string) (float64, error) {
	// Verifica se o valor contém "R$"
	if strings.Contains(value, "R$") {
		// Remove "R$" e espaços em branco
		value = strings.TrimSpace(strings.Replace(value, "R$", "", -1))
	}

	// Substitui vírgula por ponto (caso seja necessário)
	value = strings.Replace(value, ",", ".", 1)

	// Imprime o valor após a substituição para verificar o que está chegando
	fmt.Println("Valor após substituir vírgula por ponto:", value)

	// Tenta converter para float64
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.New("erro ao converter valor para float64: " + err.Error())
	}

	return result, nil
}
