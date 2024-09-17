package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Produto struct {
	Price     float64 `json:"price"`     // 8 bytes
	Promotion bool    `json:"promotion"` // 1 byte
	// 7 bytes de padding aqui
	Description string `json:"description"` // 8 bytes (ponteiro)
	Source      string `json:"source"`      // 8 bytes (ponteiro)
	Link        string `json:"link"`        // 8 bytes (ponteiro)
	ImageURL    string `json:"image_url"`   // 8 bytes (ponteiro)
}

func CrawlGoogle(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Codificar a query string
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial com timeout
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.Sleep(5*time.Second), // Pode ser ajustado se necessário
		chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	// Extrair o HTML da página
	var htmlContent string
	err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}
	log.Println("Conteúdo HTML coletado")

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	selecao := doc.Find(".tAxDx")
	if selecao.Length() == 0 {
		log.Println("Nenhum elemento com a classe .tAxDx foi encontrado.")
	} else {
		selecao.Each(func(i int, s *goquery.Selection) {
			htmlContent, err := s.Html()
			if err != nil {
				log.Printf("Erro ao obter HTML do elemento %d: %v\n", i, err)
				return
			}
			log.Printf("Elemento %d com classe tAxDx: %s\n", i, htmlContent)
		})
	}

	log.Println("Coletando o html", htmlContent)

	// Coletar dados de até 10 produtos
	var produtos []Produto
	doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
		nome := s.Find("h3.tAxDx").Text()
		link, _ := s.Find("a.xCpuod").Attr("href")
		imagemURL, _ := s.Find("img").Attr("src")
		precoStr := s.Find("span.a8Pemb").Text()
		fornecedor := "Desconhecido"

		// Verificar se o preço foi encontrado
		if precoStr == "" {
			precoStr = "0" // Valor padrão caso o preço não seja encontrado
		}

		// Formatar o preço
		precoStr = strings.TrimSpace(precoStr)
		precoStr = strings.ReplaceAll(precoStr, "R$", "")
		precoStr = strings.ReplaceAll(precoStr, ".", "")
		precoStr = strings.ReplaceAll(precoStr, ",", ".")
		precoStr = strings.ReplaceAll(precoStr, "\u00a0", "")

		// Converter preço de string para float64
		preco, err := strconv.ParseFloat(precoStr, 64)
		if err != nil {
			log.Printf("Erro ao converter preço '%s' para float64: %v", precoStr, err)
			preco = 0 // Valor padrão caso a conversão falhe
		}

		produto := Produto{
			Description: nome,
			Link:        "https://www.google.com" + link,
			ImageURL:    imagemURL,
			Price:       preco,
			Source:      fornecedor,
		}

		produtos = append(produtos, produto)
		return true
	})

	log.Println("Produtos coletados:", len(produtos))

	// Retornar a lista de produtos coletados
	return produtos, nil
}
