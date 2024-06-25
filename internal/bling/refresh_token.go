package bling

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetRefreshToken(username, password, refresh_token string) (string, error) {
	apiURL := "https://www.bling.com.br/Api/v3/oauth/token"
	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh_token},
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "1.0")
	req.SetBasicAuth(username, password)

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

	return string(bodyBytes), nil
}
