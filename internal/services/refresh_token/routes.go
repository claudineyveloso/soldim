// package refreshtoken
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
//
// 	"github.com/claudineyveloso/soldim.git/internal/bling"
// 	"github.com/gorilla/mux"
// )
//
// func RegisterRoutes(router *mux.Router) {
// 	router.HandleFunc("/refresh_token", handlePostRefreshToken).Methods(http.MethodPost)
// }
//
// func handlePostRefreshToken(w http.ResponseWriter, r *http.Request) {
// 	username := "11e56de94a8dc983459367236b79608cd941dda6"
// 	password := "26ef0f168a6c9fc7618cafacbead208a9cb4a9d2492c1f33ac4a8ccfb2c3"
// 	refreshToken := "baf38b46b469af96cceea91c91185486fcdaf92b"
//
// 	fmt.Println("Enviando requisição para gerar novo token de atualização...")
//
// 	tokenResponse, err := bling.GetRefreshToken(username, password, refreshToken)
// 	if err != nil {
// 		fmt.Println("Erro:", err)
// 		return
// 	}
//
// 	// Imprime a resposta como JSON
// 	tokenResponseJSON, err := json.MarshalIndent(tokenResponse, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao formatar JSON:", err)
// 		return
// 	}
//
// 	fmt.Println(string(tokenResponseJSON))
// }

package refreshtoken

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/refresh_token", handlePostRefreshToken).Methods(http.MethodPost)
}

func handlePostRefreshToken(w http.ResponseWriter, r *http.Request) {
	username := "11e56de94a8dc983459367236b79608cd941dda6"
	password := "26ef0f168a6c9fc7618cafacbead208a9cb4a9d2492c1f33ac4a8ccfb2c3"
	refreshToken := "0f4aad89b8abd53ab9fb18154dbc34f3e7cc8076"

	fmt.Println("Enviando requisição para gerar novo token de atualização...")

	tokenResponse, err := bling.GetRefreshToken(username, password, refreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter token: %v", err), http.StatusInternalServerError)
		return
	}

	// Converte a resposta do token para um map para facilitar a formatação JSON
	var tokenResponseMap map[string]interface{}
	if err := json.Unmarshal([]byte(tokenResponse), &tokenResponseMap); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao formatar resposta: %v", err), http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho Content-Type para application/json
	w.Header().Set("Content-Type", "application/json")

	// Envia a resposta JSON ao cliente
	if err := json.NewEncoder(w).Encode(tokenResponseMap); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao enviar resposta JSON: %v", err), http.StatusInternalServerError)
	}
}
