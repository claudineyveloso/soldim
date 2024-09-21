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
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Produto struct {
	Price       float64 `json:"price"`       // 8 bytes
	Promotion   bool    `json:"promotion"`   // 1 byte
	Description string  `json:"description"` // 8 bytes (ponteiro)
	Source      string  `json:"source"`      // 8 bytes (ponteiro)
	Link        string  `json:"link"`        // 8 bytes (ponteiro)
	ImageURL    string  `json:"image_url"`   // 8 bytes (ponteiro)
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

	log.Printf("Passou pela primeira condicao if: %s", startURL)
	for {
		// Esperar o carregamento da página
		//
		var nodes []*cdp.Node
		err = chromedp.Run(ctx, chromedp.Nodes(`div.sh-dgr__grid-result`, &nodes, chromedp.AtLeast(0)))
		if err != nil {
			log.Printf("Erro ao verificar a existência do elemento: %v", err)
			break
		}

		if len(nodes) > 0 {
			log.Printf("O elemento 'div.sh-dgr__grid-result' foi encontrado no DOM.")
		} else {
			log.Printf("O elemento 'div.sh-dgr__grid-result' NÃO foi encontrado no DOM.")
			break
		}

		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`))
		if err != nil {
			log.Printf("Erro ao esperar pela visibilidade dos resultados: %v", err)
			break
		}

		log.Printf("Passou pel segunda condicao if: %s", startURL)
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

		if doc.Find(".sh-dgr__grid-result").Length() > 0 {
			log.Printf("A classe sh-dgr__grid-result foi encontrada: %s", startURL)
		} else {
			log.Printf("A classe sh-dgr__grid-result não foi encontrada: %s", startURL)
		}

		// Extrair detalhes dos produtos
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			description := item.Find(".tAxDx").Text()
			price := formatarPreco(item.Find(".a8Pemb").Text())
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
		// Log dos produtos coletados
		log.Printf("Total de produtos coletados: %d", len(produtos))
		for _, produto := range produtos {
			log.Printf("Produto: %+v", produto)
		}
	}
	return produtos, nil
}

func formatarPreco(valor string) float64 {
	// Remover R$ e espaço não quebrável
	valor = strings.Replace(valor, "R$", "", -1)
	valor = strings.Replace(valor, "\u00a0", "", -1)

	// Substituir vírgula por ponto
	valor = strings.Replace(valor, ",", ".", -1)

	// Remover caracteres não numéricos, exceto ponto decimal
	re := regexp.MustCompile(`[^\d.]`)
	valor = re.ReplaceAllString(valor, "")

	// Converter para float64
	preco, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		log.Printf("Erro ao converter preço: %v", err)
		return 0.0
	}

	return preco
}
