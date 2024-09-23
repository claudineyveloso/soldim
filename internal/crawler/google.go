package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
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
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)

	// Criar uma nova instância do Colly com limitações
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
		colly.MaxDepth(1), // Limitar profundidade para evitar loops
	)

	// Limitar a velocidade e paralelismo
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,               // Apenas uma requisição por vez
		Delay:       2 * time.Second, // Delay entre requisições
	})

	// Slice para armazenar os produtos coletados
	var produtos []Produto

	// Tratar quando a página for visitada
	c.OnHTML("html", func(e *colly.HTMLElement) {
		produto := Produto{}

		// Coletar nome do produto
		produto.Description = e.ChildText("h4 span")
		produto.Price = 0.0 // Simplesmente atribuindo um valor fixo aqui

		// Adicionar produto ao slice
		produtos = append(produtos, produto)

		// Logar produto coletado para verificação
		log.Printf("Produto coletado: Nome: %s, Preço: %f", produto.Description, produto.Price)
	})

	// Tratar erro ao visitar a página
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Erro: %v Status Code: %d", err, r.StatusCode)
	})

	// Tratar quando a coleta for concluída
	c.OnScraped(func(r *colly.Response) {
		log.Println("Coleta finalizada:", r.Request.URL)
	})

	// Tratar HTML completo da página para salvar em arquivo
	c.OnResponse(func(r *colly.Response) {
		// Salvar HTML completo para depuração
		file, err := os.Create("pagina_completa_colly.html")
		if err != nil {
			log.Printf("Erro ao criar arquivo: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(string(r.Body))
		if err != nil {
			log.Printf("Erro ao escrever no arquivo: %v", err)
		}

		log.Println("HTML salvo em pagina_completa_colly.html")
	})

	// Iniciar a coleta visitando a página
	err := c.Visit(startURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar a página: %v", err)
	}

	// Esperar até a coleta estar finalizada
	c.Wait()

	return produtos, nil
}

// func CrawlGoogle(query string) ([]Produto, error) {
// 	// Codificar a query para ser usada na URL
// 	encodedQuery := url.QueryEscape(query)
// 	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
//
// 	// Criar uma nova instância do Colly
// 	c := colly.NewCollector(
// 		// Configurar User-Agent para evitar bloqueios
// 		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
// 	)
//
// 	// Slice para armazenar os produtos coletados
// 	var produtos []Produto
//
// 	// Tratar quando a página for visitada
// 	c.OnHTML("html", func(e *colly.HTMLElement) {
// 		produto := Produto{}
//
// 		// Coletar nome do produto
// 		produto.Description = e.ChildText("h4 span")
// 		// Coletar preço do produto (se existir)
// 		produto.Price = 0.0 // e.ChildText("span.a8Pemb")
//
// 		// Adicionar produto ao slice
// 		produtos = append(produtos, produto)
//
// 		// Logar produto coletado para verificação
// 		log.Printf("Produto coletado: Nome: %s, Preço: %f", produto.Description, produto.Price)
// 	})
//
// 	// Tratar erro ao visitar a página
// 	c.OnError(func(r *colly.Response, err error) {
// 		log.Printf("Erro: %v Status Code: %d", err, r.StatusCode)
// 	})
//
// 	// Tratar quando a coleta for concluída
// 	c.OnScraped(func(r *colly.Response) {
// 		log.Println("Coleta finalizada:", r.Request.URL)
// 	})
//
// 	// Tratar HTML completo da página para salvar em arquivo
// 	c.OnResponse(func(r *colly.Response) {
// 		// Salvar HTML completo para depuração
// 		file, err := os.Create("pagina_completa_colly.html")
// 		if err != nil {
// 			log.Printf("Erro ao criar arquivo: %v", err)
// 		}
// 		defer file.Close()
//
// 		_, err = file.WriteString(string(r.Body))
// 		if err != nil {
// 			log.Printf("Erro ao escrever no arquivo: %v", err)
// 		}
//
// 		log.Println("HTML salvo em pagina_completa_colly.html")
// 	})
//
// 	// Iniciar a coleta visitando a página
// 	err := c.Visit(startURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("falha ao visitar a página: %v", err)
// 	}
//
// 	// Esperar até a coleta estar finalizada
// 	c.Wait()
//
// 	return produtos, nil
// }

func CrawlGoogleDDD(query string) ([]Produto, error) {
	// Codificar a query para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)

	// Criar uma nova instância do Colly
	c := colly.NewCollector(
		// Configurar User-Agent para evitar bloqueios
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
	)

	// Slice para armazenar os produtos coletados
	var produtos []Produto

	// Tratar quando a página for visitada
	c.OnHTML("div.sh-dgr__content", func(e *colly.HTMLElement) {
		produto := Produto{}

		// Coletar nome do produto
		produto.Description = e.ChildText("h4 span")
		// Coletar preço do produto (se existir)
		produto.Price = 0.0 // e.ChildText("span.a8Pemb")

		// Adicionar produto ao slice
		produtos = append(produtos, produto)

		// Logar produto coletado para verificação
		log.Printf("Produto coletado: Nome: %s, Preço: %f", produto.Description, produto.Price)
	})

	// Tratar erro ao visitar a página
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Erro: %v Status Code: %d", err, r.StatusCode)
	})

	// Tratar quando a coleta for concluída
	c.OnScraped(func(r *colly.Response) {
		log.Println("Coleta finalizada:", r.Request.URL)
	})

	// Tratar HTML completo da página para salvar em arquivo
	c.OnResponse(func(r *colly.Response) {
		// Salvar HTML completo para depuração
		file, err := os.Create("pagina_completa_colly.html")
		if err != nil {
			log.Printf("Erro ao criar arquivo: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(string(r.Body))
		if err != nil {
			log.Printf("Erro ao escrever no arquivo: %v", err)
		}

		log.Println("HTML salvo em pagina_completa_colly.html")
	})

	// Iniciar a coleta visitando a página
	err := c.Visit(startURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar a página: %v", err)
	}

	// Esperar até a coleta estar finalizada
	c.Wait()

	return produtos, nil
}

func CrawlGoogleAtual(query string) ([]Produto, error) {
	// Configurar o User-Agent para simular um navegador real
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	// Aplicar o contexto com as opções configuradas
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Criar o contexto padrão a partir do allocator
	ctx, cancel = chromedp.NewContext(ctx)
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

	// Aguardar um tempo para evitar problemas com rate limiting
	time.Sleep(4 * time.Second)

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

			var link string
			if strings.HasPrefix(rawURL, "/shopping/product") {
				link = "https://www.google.com.br" + rawURL
			} else if strings.HasPrefix(rawURL, "/url?url=") {
				link = strings.TrimPrefix(rawURL, "/url?url=")
			} else {
				link = rawURL
			}

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
		time.Sleep(4 * time.Second)

		log.Printf("Total de produtos coletados: %d", len(produtos))
		for _, produto := range produtos {
			log.Printf("Produto: %+v", produto)
		}
	}

	return produtos, nil
}

func CrawlGoogleOLD(query string) ([]Produto, error) {
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

	var htmlContent string
	err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent))
	if err != nil {
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}
	log.Println("Conteúdo HTML coletado:", htmlContent)

	for {
		// Esperar o carregamento da página

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
