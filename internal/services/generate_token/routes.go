package generatetoken

import (
	"fmt"
	"log"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/generate_token", handlePostGenerateToken).Methods(http.MethodPost)
}

func handlePostGenerateToken(w http.ResponseWriter, r *http.Request) {
	clientID := "11e56de94a8dc983459367236b79608cd941dda6"
	clientSecret := "26ef0f168a6c9fc7618cafacbead208a9cb4a9d2492c1f33ac4a8ccfb2c3"
	authorizationCode := "09210d8906fe3e7c05749927482f223932e25ef3"
	tokenURL := "https://www.bling.com.br/Api/v3/oauth/token"

	token, err := bling.GetTokenFromBling(clientID, clientSecret, authorizationCode, tokenURL)
	if err != nil {
		log.Fatalf("Erro ao obter o token: %v", err)
	}

	fmt.Printf("Token de acesso: %s\n", token)
}
