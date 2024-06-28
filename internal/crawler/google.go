package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Produto struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Source      string  `json:"source"`
	Link        string  `json:"link"`
	ImageURL    string  `json:"image_url"`
	Promotion   bool    `json:"promotion"`
}

func CrawlGoogle(query string) ([]Produto, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var produtos []Produto

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Printf("Iniciando visita: %s", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
	if err != nil {
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	for {
		// Esperar o carregamento da página
		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`))
		if err != nil {
			log.Printf("Erro ao esperar pela visibilidade dos resultados: %v", err)
			break
		}

		// Extrair o HTML da página
		var htmlContent string
		err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent))
		if err != nil {
			return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
		}

		// Parsear o HTML com goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
		}

		// Extrair detalhes dos produtos
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			description := item.Find(".tAxDx").Text()
			priceText := item.Find(".a8Pemb").Text()
			log.Printf("Raw price text: %s", priceText)

			price, err := formatarPreco(priceText)
			if err != nil {
				log.Printf("Erro ao formatar o preço: %v", err)
				price = 0.0
			}
			log.Printf("Formatted price: %f", price)

			rawURL, _ := item.Find("a").Attr("href")
			imageURL, _ := item.Find(".ArOc1c img").Attr("src")
			promotionText := strings.TrimSpace(item.Find(".fAcMNb span.Ib8pOd").Text())

			source := ""
			item.Find(".aULzUe").Contents().Each(func(i int, s *goquery.Selection) {
				if goquery.NodeName(s) != "style" {
					source = strings.TrimSpace(s.Text())
				}
			})

			// Processar a URL conforme a lógica solicitada
			var link string
			if strings.HasPrefix(rawURL, "/shopping/product") {
				link = "https://www.google.com.br" + rawURL
			} else if strings.HasPrefix(rawURL, "/url?url=") {
				link = strings.TrimPrefix(rawURL, "/url?url=")
			} else {
				link = rawURL
			}

			// Verificar se o texto da promoção é "PROMOÇÃO"
			promotion := promotionText == "PROMOÇÃO"

			produto := Produto{
				Description: strings.TrimSpace(description),
				Price:       price,
				Source:      source,
				Link:        link,
				ImageURL:    imageURL,
				Promotion:   promotion,
			}
			produtos = append(produtos, produto)
			log.Printf("Produto encontrado: %+v\n", produto)
		})

		// Verificar se há uma próxima página
		var nextPageExists bool
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('a#pnnext') !== null`, &nextPageExists))
		if err != nil {
			return nil, fmt.Errorf("falha ao verificar a próxima página: %v", err)
		}

		if !nextPageExists {
			break
		}

		// Navegar para a próxima página
		err = chromedp.Run(ctx, chromedp.Click(`a#pnnext`, chromedp.NodeVisible))
		if err != nil {
			return nil, fmt.Errorf("falha ao navegar para a próxima página: %v", err)
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		time.Sleep(2 * time.Second)
	}

	// Log dos produtos coletados
	log.Printf("Total de produtos coletados: %d", len(produtos))
	for _, produto := range produtos {
		log.Printf("Produto: %+v", produto)
	}

	return produtos, nil
}

func formatarPreco(valor string) (float64, error) {
	// Remover R$ e espaço não quebrável
	log.Printf("Raw valor: %s", valor)
	valor = strings.Replace(valor, "R$", "", -1)
	valor = strings.Replace(valor, "\u00a0", "", -1)
	log.Printf("Valor after removing R$ and non-breaking space: %s", valor)

	// Substituir vírgula por ponto
	valor = strings.Replace(valor, ".", "", -1)  // Remove thousands separator
	valor = strings.Replace(valor, ",", ".", -1) // Replace decimal comma with dot
	log.Printf("Valor after replacing comma with dot: %s", valor)

	// Remover caracteres não numéricos, exceto ponto decimal
	re := regexp.MustCompile(`[^\d.]`)
	valor = re.ReplaceAllString(valor, "")
	log.Printf("Valor after removing non-numeric characters: %s", valor)

	// Converter para float64
	preco, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
	}
	return preco, err
}
