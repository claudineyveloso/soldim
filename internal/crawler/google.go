package crawler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Produto struct {
	Nome   string
	Valor  string
	Fonte  string
	URL    string
	Imagem string
}

func CrawlGoogle(query string) []Produto {
	var produtos []Produto
	var produtoCount int

	// Criar um novo coletor para a pesquisa no Google
	c := colly.NewCollector(
		colly.AllowedDomains("www.google.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, como Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Criar um coletor secundário para seguir o link de compras e lidar com a paginação
	shoppingCollector := colly.NewCollector(
		colly.AllowedDomains("www.google.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, como Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Manipulador para quando uma página HTML é visitada no Google
	c.OnHTML("a", func(e *colly.HTMLElement) {
		linkText := e.Text
		linkHref := e.Attr("href")

		// Verificar se o link é para a seção de compras
		if strings.Contains(linkText, "Shopping") {
			fullURL := "https://www.google.com" + linkHref
			log.Printf("Seguindo o link de compras: %s\n", fullURL)
			shoppingCollector.Visit(fullURL)
		}
	})

	// Manipulador para quando uma página HTML é visitada na seção de compras
	shoppingCollector.OnHTML("div.sh-dgr__content", func(e *colly.HTMLElement) {
		nome := e.ChildText(".EI11Pd h3.tAxDx")
		valor := e.ChildText(".a8Pemb")
		fonte := e.ChildText(".aULzUe.IuHnof")
		url := e.ChildAttr("a", "href")
		imagem := e.ChildAttr(".ArOc1c img", "src")

		produto := Produto{
			Nome:   nome,
			Valor:  valor,
			Fonte:  fonte,
			URL:    url,
			Imagem: imagem,
		}
		produtos = append(produtos, produto)
		produtoCount++
		log.Printf("Produto encontrado: %+v\n", produto)
	})

	// Manipulador para quando uma requisição falha
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Algo deu errado:", err)
	})

	shoppingCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Algo deu errado no coletor de compras:", err)
	})

	// Manipulador para lidar com a paginação
	shoppingCollector.OnHTML("a#pnnext", func(e *colly.HTMLElement) {
		linkHref := e.Attr("href")
		if linkHref != "" {
			fullURL := "https://www.google.com" + linkHref
			time.Sleep(2 * time.Second) // Adicione um pequeno atraso para evitar problemas com rate limiting
			log.Printf("Seguindo para a próxima página: %s\n", fullURL)
			shoppingCollector.Visit(fullURL)
		}
	})

	// Iniciar a pesquisa no Google
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s", query)
	err := c.Visit(startURL)
	if err != nil {
		log.Fatalf("Falha ao iniciar a visita: %v", err)
	}

	// Esperar até que todas as visitas sejam concluídas
	c.Wait()
	shoppingCollector.Wait()

	// // Exibir os resultados coletados
	// fmt.Printf("Total de produtos encontrados: %d\n", produtoCount)
	// for _, produto := range produtos {
	// 	fmt.Printf("Nome: %s\n", produto.Nome)
	// 	fmt.Printf("Valor: %s\n", produto.Valor)
	// 	fmt.Printf("Fonte: %s\n", produto.Fonte)
	// 	fmt.Printf("URL: %s\n", produto.URL)
	// 	fmt.Printf("Imagem: %s\n", produto.Imagem)
	// 	fmt.Println("-------------------------------")
	// }

	return produtos
}
