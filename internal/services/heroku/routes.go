package heroku

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var token = os.Getenv("ACCESS_TOKEN_BLING")

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api_contacts_bling", handleGetAPIHeroku).Methods(http.MethodGet)
}

func handleGetAPIHeroku(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("API Acesso aos Contatos da Bling!"))
	if err != nil {
		http.Error(w, "Erro ao escrever resposta", http.StatusInternalServerError)
		return
	}
}

// func handleGetAPIHeroku(w http.ResponseWriter, r *http.Request) {
// 	contacts, err := bling.GetContactsFromBling(token)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Erro ao obter contatos: %v", err), http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Definir o tipo de conteúdo da resposta como JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
//
// 	// Converter []types.Contact para JSON
// 	jsonData, err := json.Marshal(contacts)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Erro ao converter contatos para JSON: %v", err), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.Write(jsonData) // Enviar a resposta JSON
// }
//
// func GetContactsFromBling(bearerToken string) ([]types.Contact, error) {
// 	req, err := http.NewRequest("GET", "https://bling.com.br/Api/v3/contatos", nil)
// 	if err != nil {
// 		fmt.Println("Erro ao criar requisição:", err)
// 		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
// 	}
// 	req.Header.Set("Authorization", "Bearer "+bearerToken)
//
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Erro ao enviar requisição:", err)
// 		return nil, fmt.Errorf("erro ao enviar requisição: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode == http.StatusUnauthorized {
// 		log.Println("Token expirado. Tentando renovar o token...")
// 		newToken, err := RefreshToken()
// 		if err != nil {
// 			return nil, fmt.Errorf("erro ao renovar token: %v", err)
// 		}
// 		return GetContactsFromBling(newToken)
// 	}
//
// 	if resp.StatusCode != http.StatusOK {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		bodyString := string(bodyBytes)
// 		log.Printf("Status Code: %d", resp.StatusCode)
// 		log.Printf("Response Body: %s", bodyString)
// 		return nil, fmt.Errorf("falha na requisição: %s", bodyString)
// 	}
//
// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
// 	}
//
// 	var responseData types.ContactResponse
// 	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
// 		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
// 	}
//
// 	log.Printf("Número de contatos retornados: %d\n", len(responseData.Data))
//
// 	return responseData.Data, nil
// }
//
// func RefreshToken() (string, error) {
// 	username := "11e56de94a8dc983459367236b79608cd941dda6"
// 	password := "26ef0f168a6c9fc7618cafacbead208a9cb4a9d2492c1f33ac4a8ccfb2c3"
// 	refreshToken := "01f24f4049658fe4398ec63102715e860ef632f6"
//
// 	tokenResponse, err := bling.GetRefreshToken(username, password, refreshToken)
// 	if err != nil {
// 		return "", fmt.Errorf("erro ao obter token: %v", err)
// 	}
//
// 	// Extraia o novo token da resposta
// 	var tokenResponseMap map[string]interface{}
// 	if err := json.Unmarshal([]byte(tokenResponse), &tokenResponseMap); err != nil {
// 		return "", fmt.Errorf("erro ao formatar resposta: %v", err)
// 	}
//
// 	// Aqui você precisará ajustar o código para extrair o token do map corretamente
// 	newToken := tokenResponseMap["access_token"].(string) // Certifique-se que o campo é esse mesmo
//
// 	// Atualize o token global ou variável de ambiente
// 	return newToken, nil
// }
