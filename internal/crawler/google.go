package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
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

	// Definir contexto com timeout de 15 minutos
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
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

	// Extrair o HTML da página
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	// Logar o conteúdo HTML
	// log.Println("Conteúdo HTML coletado:", htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Localizar e atribuir texto das tags à variável nome
	var description, price, source, image string

	doc.Find("span.pymv4e, h3.tAxDx").Each(func(i int, s *goquery.Selection) {
		description += s.Text() + " "
	})

	doc.Find("span.lmQWe, span.a8Pemb").Each(func(i int, s *goquery.Selection) {
		price += s.Text() + " "
	})

	// Extraindo fontes (local e Heroku)
	doc.Find("span.zPEcBd, .aULzUe.IuHnof").Each(func(i int, s *goquery.Selection) {
		source += s.Text() + " "
	})

	doc.Find(".D6nsM, .ArOc1c").Each(func(i int, s *goquery.Selection) {
		if imgSrc, exists := s.Find("img").Attr("src"); exists {
			image += imgSrc + " "
		} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
			image += imgDataSrc + " "
		}
	})

	// Logar o nome extraído
	log.Println("Descrição extraída:", description)
	log.Println("Preço extraído:", price)
	log.Println("Fonte extraída:", source)
	log.Println("Imagens extraídas:", image)
	// Retornar nil, já que não estamos processando os produtos por enquanto
	return nil, nil
}

func CrawlGoogleBBB(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
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

	// Extrair o HTML da página
	var htmlContent string
	err := chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}
	log.Println("Conteúdo HTML coletado")

	log.Println("Coletando o html", htmlContent)

	// Parsear o HTML usando goquery
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	// if err != nil {
	// 	log.Println("Falha ao parsear HTML:", err)
	// 	return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	// }

	// selecao := doc.Find(".tAxDx")
	// if selecao.Length() == 0 {
	// 	log.Println("Nenhum elemento com a classe .tAxDx foi encontrado.")
	// } else {
	// 	selecao.Each(func(i int, s *goquery.Selection) {
	// 		htmlContent, err := s.Html()
	// 		if err != nil {
	// 			log.Printf("Erro ao obter HTML do elemento %d: %v\n", i, err)
	// 			return
	// 		}
	// 		log.Printf("Elemento %d com classe tAxDx: %s\n", i, htmlContent)
	// 	})
	// }
	//
	// // Coletar dados de até 10 produtos
	// var produtos []Produto
	// doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	nome := s.Find("h3.tAxDx").Text()
	// 	link, _ := s.Find("a.xCpuod").Attr("href")
	// 	imagemURL, _ := s.Find("img").Attr("src")
	// 	precoStr := s.Find("span.a8Pemb").Text()
	// 	fornecedor := "Desconhecido"
	//
	// 	// Verificar se o preço foi encontrado
	// 	if precoStr == "" {
	// 		precoStr = "0" // Valor padrão caso o preço não seja encontrado
	// 	}
	//
	// 	// Formatar o preço
	// 	precoStr = strings.TrimSpace(precoStr)
	// 	precoStr = strings.ReplaceAll(precoStr, "R$", "")
	// 	precoStr = strings.ReplaceAll(precoStr, ".", "")
	// 	precoStr = strings.ReplaceAll(precoStr, ",", ".")
	// 	precoStr = strings.ReplaceAll(precoStr, "\u00a0", "")
	//
	// 	// Converter preço de string para float64
	// 	preco, err := strconv.ParseFloat(precoStr, 64)
	// 	if err != nil {
	// 		log.Printf("Erro ao converter preço '%s' para float64: %v", precoStr, err)
	// 		preco = 0 // Valor padrão caso a conversão falhe
	// 	}
	//
	// 	produto := Produto{
	// 		Description: nome,
	// 		Link:        "https://www.google.com" + link,
	// 		ImageURL:    imagemURL,
	// 		Price:       preco,
	// 		Source:      fornecedor,
	// 	}
	//
	// 	produtos = append(produtos, produto)
	//	return true
	//})
	//

	// log.Println("Produtos coletados:", len(produtos))

	// Retornar a lista de produtos coletados
	// return produtos, nil
	return nil, nil
}
