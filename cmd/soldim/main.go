package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/claudineyveloso/soldim.git/cmd/api"
	"github.com/claudineyveloso/soldim.git/cmd/db"
	"github.com/claudineyveloso/soldim.git/internal/configs"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Soldim"),
		newrelic.ConfigLicense("e507d764f5d69169cca78f92c9cdbe43FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Fatalln("Erro ao iniciar New Relic:", err)
	}

	// Verifica se o app está conectado ao New Relic
	if app == nil {
		log.Fatalln("New Relic app não foi inicializado corretamente.")
	}

	// Função para processar requisições, agora monitoradas pelo New Relic
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/soldim-newrelic", handlerFunction))

	// Carrega as configurações do ambiente
	cfg := configs.Envs

	// Conecta ao banco de dados
	dbConn, err := db.NewPostgresSQLStorage(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Inicializa o banco de dados (verificação de conexão)
	initStorage(dbConn)

	// Inicializa o servidor da API na porta correta
	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), dbConn)
	if err := server.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar a conexão com o banco de dados: %v", err)
	}
	log.Println("DB: Conexão estabelecida com sucesso!")
}

func handlerFunction(w http.ResponseWriter, r *http.Request) {
	// Lógica da requisição
	w.Write([]byte("New Relic está monitorando!"))
}
