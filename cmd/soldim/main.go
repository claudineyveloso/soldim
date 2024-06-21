package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claudineyveloso/soldim.git/cmd/api"
	"github.com/claudineyveloso/soldim.git/cmd/db"
	"github.com/claudineyveloso/soldim.git/internal/configs"
)

func main() {
	cfg := configs.Config{
		PublicHost: configs.Envs.PublicHost,
		Port:       configs.Envs.Port,
		DBUser:     configs.Envs.DBUser,
		DBPassword: configs.Envs.DBPassword,
		DBName:     configs.Envs.DBName,
	}
	db, err := db.NewPostgresSQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")
}

// query := "Capacete Protork"
// produtos := crawler.CrawlGoogle(query)
//
// fmt.Printf("Total de produtos encontrados: %d\n", len(produtos))
// for _, produto := range produtos {
// 	fmt.Printf("Nome: %s\n", produto.Nome)
// 	fmt.Printf("Valor: %s\n", produto.Valor)
// 	fmt.Printf("Fonte: %s\n", produto.Fonte)
// 	fmt.Printf("URL: %s\n", produto.URL)
// 	fmt.Printf("Imagem: %s\n", produto.Imagem)
// 	fmt.Println("-------------------------------")
// }
//
// fmt.Printf("Total de produtos encontrados: %d\n", len(produtos))
// bearerToken := "cafdf8d595a1cd9524c8fc7dc1b1e19c676fc8da"
// produtos, err := crawler.GetProductsFromBling(bearerToken)
// if err != nil {
// 	log.Fatalf("Erro ao obter produtos: %v", err)
// }
//
// if len(produtos) == 0 {
// 	log.Println("Nenhum produto encontrado.")
// } else {
// 	for _, produto := range produtos {
// 		preco := fmt.Sprintf("%v", produto.Preco) // Converter interface{} para string
//
// 		fmt.Printf("ID: %d\n", produto.ID)
// 		fmt.Printf("Nome: %s\n", produto.Nome)
// 		fmt.Printf("Código: %s\n", produto.Codigo)
// 		fmt.Printf("Preço: %s\n", preco)
// 		fmt.Printf("Tipo: %s\n", produto.Tipo)
// 		fmt.Printf("Situação: %s\n", produto.Situacao)
// 		fmt.Printf("Formato: %s\n", produto.Formato)
// 		fmt.Printf("Descrição Curta: %s\n", produto.DescricaoCurta)
// 		fmt.Printf("Imagem URL: %s\n", produto.ImagemURL)
// 		fmt.Println("------------------------------")
// 	}
// }
//}
