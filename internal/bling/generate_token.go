package bling

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	AccessToken  string `json:"id"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// GetTokenFromBling gera o token de acesso usando client_id, client_secret, authorization_code
func GetTokenFromBling(clientID, clientSecret, authorizationCode, tokenURL string) (string, error) {
	// Codifica o clientID e clientSecret em Base64 para a autorização Basic
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Monta os dados do corpo da requisição
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authorizationCode)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return "", fmt.Errorf("falha na requisição: %s", bodyString)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token não encontrado na resposta")
	}

	return accessToken, nil
}
