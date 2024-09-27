package crawler

import (
	"context"
	"encoding/base64"
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

	"github.com/gocolly/colly/v2"
)

type Produto struct {
	Price       float64 `json:"price"`       // 8 bytes
	Promotion   bool    `json:"promotion"`   // 1 byte
	Description string  `json:"description"` // 8 bytes (ponteiro)
	Source      string  `json:"source"`      // 8 bytes (ponteiro)
	Link        string  `json:"link"`        // 8 bytes (ponteiro)
	ImageURL    string  `json:"image_url"`   // 8 bytes (ponteiro)
}

var links []string

func CrawlGoogle(query string) ([]Produto, error) {
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)

	var produtos []Produto
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
		colly.MaxDepth(2), // Limitar profundidade para evitar loops

	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	})

	// Limitar a velocidade e paralelismo
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,               // Apenas uma requisição por vez
		Delay:       5 * time.Second, // Delay entre requisições
	})

	c.OnHTML(".sh-dgr__grid-result", func(e *colly.HTMLElement) {
		// elementoEncontrado = true

		// Itera sobre os produtos dentro da classe principal
		e.ForEach(".i0X6df", func(_ int, el *colly.HTMLElement) {
			description := el.ChildText("h3.tAxDx")
			spans := e.DOM.Find("span.a8Pemb")
			price := formatarPreco(spans.Eq(0).Text())
			source := e.DOM.Find(".aULzUe.IuHnof")
			link := el.ChildAttr("a", "href")
			imageURL := el.ChildAttr(".ArOc1c img", "src")
			promotion := false

			// // Exibe o nome do produto no terminal
			// fmt.Printf("Product Name: %s\n", description)
			// fmt.Printf("Product Price: %f\n", price)
			// fmt.Printf("Product Source: %s\n", extrairSource(source.Eq(0).Text()))
			// fmt.Printf("Product Image: %s\n", imageURL)
			// fmt.Printf("Product URL: %s\n", link)
			// fmt.Printf("Product Img: %s\n", img)
			// fmt.Printf("Product Promotion: %v\n", promotion)
			//
			produto := Produto{
				Description: strings.TrimSpace(description),
				Price:       price,
				Source:      extrairSource(source.Text()),
				Link:        link,
				ImageURL:    imageURL,
				Promotion:   promotion,
			}
			produtos = append(produtos, produto)
		})
		// log.Println("Elemento '.sh-dgr__grid-result' encontrado!")
	})
	c.OnHTML("a#pnnext", func(e *colly.HTMLElement) {
		log.Println("Paginacao encontrada")
		time.Sleep(4 * time.Second)
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Navegando para a próxima página:", nextPage)
		err := c.Visit(nextPage)
		if err != nil {
			log.Println("Erro ao visitar a próxima página:", err)
		}
	})

	// // Tratar erro ao visitar a página
	// c.OnError(func(_ *colly.Response, err error) {
	// 	// log.Println("Erro ao coletar links:", err)
	// })
	//
	// // Tratar quando a coleta for concluída
	// c.OnScraped(func(_ *colly.Response) {
	// 	if !elementoEncontrado {
	// 		fmt.Println("Elemento '.sh-dgr__grid-result' não encontrado!")
	// 	}
	// 	if !outroElemento {
	// 		fmt.Println("Elemento '.acCJ4b' não encontrado!")
	// 	}
	// 	if !eleKZmu8e {
	// 		fmt.Println("Elemento '.KZmu8e' não encontrado!")
	// 	}
	// })
	//
	// Iniciar a coleta visitando a página
	err := c.Visit(startURL)
	if err != nil {
		log.Fatalf("Erro ao visitar a página: %v", err)
	}
	log.Println("Visitando a url", startURL)

	c.Wait()

	// Codificar a query string para ser usada na URL
	return produtos, nil
}

func CrawlGoogleAtualFuncionando(query string) ([]Produto, error) {
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)

	// Criar uma nova instância do Colly com limitações
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
		colly.MaxDepth(2), // Limitar profundidade para evitar loops
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	})

	// Limitar a velocidade e paralelismo
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,               // Apenas uma requisição por vez
		Delay:       5 * time.Second, // Delay entre requisições
	})

	// Slice para armazenar os produtos coletados
	var produtos []Produto

	// Tratar quando a página for visitada
	c.OnHTML("div.Ez5pwe", func(e *colly.HTMLElement) {
		produto := Produto{}
		precoStr := e.ChildText("span.lmQWe")
		if precoStr == "" {
			precoStr = e.ChildText("span.lmQWe.YQkzwf")
		}
		// precoStr := e.ChildText("span.lmQWe.YQkzwf.pVBUqb")
		precoStr = strings.ReplaceAll(precoStr, "$", "")  // Remover o símbolo da moeda, se necessário
		precoStr = strings.ReplaceAll(precoStr, "R$", "") // Remover o símbolo da moeda, se necessário
		precoStr = strings.ReplaceAll(precoStr, ".", "")  // Remover pontos, se o formato for R$ 1.234,56
		precoStr = strings.ReplaceAll(precoStr, ",", ".")
		preco, err := strconv.ParseFloat(precoStr, 64)
		if err != nil {
			log.Printf("Erro ao converter preço para float: %v", err)
			produto.Price = 0.0 // Atribuir valor padrão em caso de erro
		} else {
			produto.Price = preco
		}
		// Coletar nome do produto
		produto.Description = e.ChildText("div.gkQHve")
		// produto.Description = e.ChildText("div.gkQHve.SsM98d.RmEs5b")
		produto.Source = e.ChildText("span.WJMUdc")

		e.ForEach("div.JK3kIe img", func(_ int, imgElement *colly.HTMLElement) {
			imageSrc := imgElement.Attr("src")

			if strings.HasPrefix(imageSrc, "data:image/") {
				// A imagem está em base64
				log.Println("Imagem base64 encontrada:", imageSrc)

				// Separar a metadata (data:image/webp;base64,) do código base64
				data := strings.Split(imageSrc, ",")[1]

				// Decodificar o base64
				decodedImage, err := base64.StdEncoding.DecodeString(data)
				if err != nil {
					log.Println("Erro ao decodificar imagem:", err)
				} else {
					// Salvar o arquivo como imagem, por exemplo, como "imagem.webp"
					err = os.WriteFile("imagem.webp", decodedImage, 0o644)
					if err != nil {
						log.Println("Erro ao salvar imagem:", err)
					} else {
						log.Println("Imagem base64 salva com sucesso.")
					}
				}
			} else {
				// Caso seja uma URL normal, processar normalmente
				log.Println("URL da imagem:", imageSrc)
				produto.ImageURL = imageSrc
			}
		})

		// Adicionar produto ao slice
		produtos = append(produtos, produto)

		// Logar produto coletado para verificação
		log.Printf("Produto coletado: Nome: %s, Preço: %f, Fonte: %s, Imagem: %s", produto.Description, produto.Price, produto.Source, produto.ImageURL)
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
			return
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

func formatarPreco_OLD(valor string) float64 {
	// Remover "R$" e espaços
	valor = strings.TrimSpace(strings.Replace(valor, "R$", "", -1))

	// Substituir vírgula por ponto
	valor = strings.Replace(valor, ",", ".", -1)

	// Remover qualquer caractere que não seja número ou ponto
	re := regexp.MustCompile(`[^\d.]`)
	valor = re.ReplaceAllString(valor, "")

	// Verifica se o valor resultante é válido antes de tentar convertê-lo
	if valor == "" {
		log.Println("Valor vazio após formatação")
		return 0.0
	}

	// Converter o valor formatado para float64
	preco, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		log.Printf("Erro ao converter preço '%s' para float: %v", valor, err)
		return 0.0
	}

	return preco
}

func extrairSource(input string) string {
	// Procurar pelo último fechamento de chave "}"
	index := strings.LastIndex(input, "}")
	if index == -1 {
		// Se não houver "}", retorna a string original
		return input
	}

	// Extrair o conteúdo após a última chave "}"
	result := input[index+1:]

	// Remover espaços extras e caracteres desnecessários
	result = strings.TrimSpace(result)

	// Retorna o valor limpo
	return result
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
