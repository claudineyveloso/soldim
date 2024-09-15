package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

func CrawlGoogleDDD(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", true),
		chromedp.Flag("remote-debugging-port", "9222"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var produtos []Produto
	produtoChan := make(chan Produto)
	var wg sync.WaitGroup

	// Número fixo de workers para processar produtos
	const numWorkers = 5

	// Função para processar produtos
	processarProduto := func(item *goquery.Selection) {
		defer wg.Done()
		description := item.Find(".tAxDx").Text()
		priceText := item.Find(".a8Pemb").Text()
		log.Println("Raw price text:", priceText)

		price, err := formatarPreco(priceText)
		if err != nil {
			log.Println("Erro ao formatar o preço:", err)
			price = 0.0
		}
		log.Println("Formatted price:", price)

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
		produtoChan <- produto
		log.Println("Produto encontrado:", produto)
	}

	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for produto := range produtoChan {
				// Aqui você pode usar o produto extraído, não `*goquery.Selection`
				// Isso pode ser processado ou usado conforme necessário
				produtos = append(produtos, produto)
			}
		}()
	}

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		close(produtoChan) // Fechar o canal ao encontrar um erro
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	for {
		// Esperar o carregamento da página com timeout específico
		log.Println("Esperando os resultados da página da coleta...")
		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade dos resultados:", err)

			// Adicionando log para verificar o estado da página
			log.Println("Verificando o HTML da página após falha...")
			var htmlContent string
			err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
			if err != nil {
				log.Println("Erro ao extrair HTML:", err)
			} else {
				log.Println("HTML da página:", htmlContent) // Adicionando log para diagnóstico
			}

			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao esperar pela visibilidade dos resultados: %v", err)
		}

		// Extrair o HTML da página
		log.Println("Extraindo HTML da página...")
		var htmlContent string
		err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
		if err != nil {
			log.Println("Falha ao extrair HTML:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
		}

		// Log do tamanho do HTML extraído para verificar se está completo
		log.Println("Tamanho do HTML extraído:", len(htmlContent), "bytes")

		log.Println("HTML extraído com sucesso. Processando o HTML...")

		// Parsear o HTML com goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			log.Println("Falha ao parsear HTML:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
		}

		// Extrair detalhes dos produtos
		log.Println("Extraindo produtos da página...")
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			wg.Add(1)
			go processarProduto(item) // Passar o *goquery.Selection para processarProduto
		})

		// Verificar se há uma próxima página com timeout específico
		log.Println("Verificando se há uma próxima página...")
		var nextPageExists bool

		// Verificando se o botão de próxima página existe
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('a#pnnext') !== null`, &nextPageExists))
		if err != nil {
			log.Println("Erro ao verificar próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao verificar próxima página: %v", err)
		}

		if !nextPageExists {
			log.Println("Não há mais páginas para navegar.")
			break
		}

		log.Println("Próxima página encontrada, tentando navegar...")

		// Verificando se o elemento está visível
		err = chromedp.Run(ctx, chromedp.WaitVisible(`a#pnnext`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade do botão de próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao esperar pela visibilidade do botão de próxima página: %v", err)
		}

		// Navegar para a próxima página
		err = chromedp.Run(ctx, chromedp.Click(`a#pnnext`, chromedp.ByQuery, chromedp.NodeVisible))
		if err != nil {
			log.Println("Falha ao navegar para a próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao navegar para a próxima página: %v", err)
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		log.Println("Aguardando para evitar rate limiting...")
		time.Sleep(1 * time.Second)
	}

	// Fechar o canal dos produtos após o WaitGroup terminar
	go func() {
		wg.Wait()
		close(produtoChan)
	}()

	// Coletar todos os produtos do canal
	for produto := range produtoChan {
		produtos = append(produtos, produto)
	}

	// Log dos produtos coletados
	log.Println("Total de produtos coletados:", len(produtos))
	for _, produto := range produtos {
		log.Println("Produto:", produto)
	}

	return produtos, nil
}

func CrawlGoogleCCC(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true), // Adicione outras flags necessárias aqui
		chromedp.Flag("headless", true),
		chromedp.Flag("remote-debugging-port", "9222"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var produtos []Produto
	produtoChan := make(chan Produto)
	var wg sync.WaitGroup

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	for {
		// Esperar o carregamento da página com timeout específico
		log.Println("Esperando os resultados da página da coleta...")
		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade dos resultados:", err)

			// Adicionando log para verificar o estado da página
			log.Println("Verificando o HTML da página após falha...")
			var htmlContent string
			err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
			if err != nil {
				log.Println("Erro ao extrair HTML:", err)
			} else {
				log.Println("HTML da página:", htmlContent) // Adicionando log para diagnóstico
			}

			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao esperar pela visibilidade dos resultados: %v", err)
		}

		// Extrair o HTML da página
		log.Println("Extraindo HTML da página...")
		var htmlContent string
		err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
		if err != nil {
			log.Println("Falha ao extrair HTML:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
		}

		// Log do tamanho do HTML extraído para verificar se está completo
		log.Println("Tamanho do HTML extraído:", len(htmlContent), "bytes")

		log.Println("HTML extraído com sucesso. Processando o HTML...")

		// Parsear o HTML com goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			log.Println("Falha ao parsear HTML:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
		}

		// Extrair detalhes dos produtos
		log.Println("Extraindo produtos da página...")
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			wg.Add(1)
			go func(item *goquery.Selection) {
				defer wg.Done()
				description := item.Find(".tAxDx").Text()
				priceText := item.Find(".a8Pemb").Text()
				log.Println("Raw price text:", priceText)

				price, err := formatarPreco(priceText)
				if err != nil {
					log.Println("Erro ao formatar o preço:", err)
					price = 0.0
				}
				log.Println("Formatted price:", price)

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
				produtoChan <- produto
				log.Println("Produto encontrado:", produto)
			}(item)
		})

		// Verificar se há uma próxima página com timeout específico
		log.Println("Verificando se há uma próxima página...")
		var nextPageExists bool

		// Verificando se o botão de próxima página existe
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('a#pnnext') !== null`, &nextPageExists))
		if err != nil {
			log.Println("Erro ao verificar próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao verificar próxima página: %v", err)
		}

		if !nextPageExists {
			log.Println("Não há mais páginas para navegar.")
			break
		}

		log.Println("Próxima página encontrada, tentando navegar...")

		// Verificando se o elemento está visível
		err = chromedp.Run(ctx, chromedp.WaitVisible(`a#pnnext`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade do botão de próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("erro ao esperar pela visibilidade do botão de próxima página: %v", err)
		}

		// Navegar para a próxima página
		err = chromedp.Run(ctx, chromedp.Click(`a#pnnext`, chromedp.ByQuery, chromedp.NodeVisible))
		if err != nil {
			log.Println("Falha ao navegar para a próxima página:", err)
			close(produtoChan) // Fechar o canal ao encontrar um erro
			return nil, fmt.Errorf("falha ao navegar para a próxima página: %v", err)
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		log.Println("Aguardando para evitar rate limiting...")
		time.Sleep(1 * time.Second)
	}

	// Fechar o canal dos produtos após o WaitGroup terminar
	go func() {
		wg.Wait()
		close(produtoChan)
	}()

	// Coletar todos os produtos do canal
	for produto := range produtoChan {
		produtos = append(produtos, produto)
	}

	// Log dos produtos coletados
	log.Println("Total de produtos coletados:", len(produtos))
	for _, produto := range produtos {
		log.Println("Produto:", produto)
	}

	return produtos, nil
}

func CrawlGoogleOLD(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true), // Adicione outras flags necessárias aqui
		chromedp.Flag("headless", true),
		chromedp.Flag("remote-debugging-port", "9222"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var produtos []Produto

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// err := chromedp.Run(ctx,
	// 	chromedp.ActionFunc(func(ctx context.Context) error {
	// 		return chromedp.Navigate(startURL).Do(ctx)
	// 	}),
	// )
	// if err != nil {
	// 	log.Println("Falha ao iniciar a visita:", err)
	// 	return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	// }

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	for {
		// Esperar o carregamento da página com timeout específico
		log.Println("Esperando os resultados da página da coleta...")
		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade dos resultados:", err)

			// Adicionando log para verificar o estado da página
			log.Println("Verificando o HTML da página após falha...")
			var htmlContent string
			err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
			if err != nil {
				log.Println("Erro ao extrair HTML:", err)
			} else {
				log.Println("HTML da página:", htmlContent) // Adicionando log para diagnóstico
			}

			break
		}

		// Extrair o HTML da página
		log.Println("Extraindo HTML da página...")
		var htmlContent string
		err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
		if err != nil {
			log.Println("Falha ao extrair HTML:", err)
			return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
		}

		// Log do tamanho do HTML extraído para verificar se está completo
		log.Println("Tamanho do HTML extraído:", len(htmlContent), "bytes")

		log.Println("HTML extraído com sucesso. Processando o HTML...")

		// Parsear o HTML com goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			log.Println("Falha ao parsear HTML:", err)
			return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
		}

		// Extrair detalhes dos produtos
		log.Println("Extraindo produtos da página...")
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			description := item.Find(".tAxDx").Text()
			priceText := item.Find(".a8Pemb").Text()
			log.Println("Raw price text:", priceText)

			price, err := formatarPreco(priceText)
			if err != nil {
				log.Println("Erro ao formatar o preço:", err)
				price = 0.0
			}
			log.Println("Formatted price:", price)

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
			log.Println("Produto encontrado:", produto)
		})

		// Verificar se há uma próxima página com timeout específico
		log.Println("Verificando se há uma próxima página...")
		var nextPageExists bool

		// Verificando se o botão de próxima página existe
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('a#pnnext') !== null`, &nextPageExists))
		if err != nil {
			log.Println("Erro ao verificar próxima página:", err)
			return nil, fmt.Errorf("erro ao verificar próxima página: %v", err)
		}

		if !nextPageExists {
			log.Println("Não há mais páginas para navegar.")
			break
		}

		log.Println("Próxima página encontrada, tentando navegar...")

		// Verificando se o elemento está visível
		err = chromedp.Run(ctx, chromedp.WaitVisible(`a#pnnext`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade do botão de próxima página:", err)
			return nil, fmt.Errorf("erro ao esperar pela visibilidade do botão de próxima página: %v", err)
		}

		// Navegar para a próxima página
		err = chromedp.Run(ctx, chromedp.Click(`a#pnnext`, chromedp.ByQuery, chromedp.NodeVisible))
		if err != nil {
			log.Println("Falha ao navegar para a próxima página:", err)
			return nil, fmt.Errorf("falha ao navegar para a próxima página: %v", err)
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		log.Println("Aguardando para evitar rate limiting...")
		time.Sleep(3 * time.Second)
	}

	// Log dos produtos coletados
	log.Println("Total de produtos coletados:", len(produtos))
	for _, produto := range produtos {
		log.Println("Produto:", produto)
	}

	return produtos, nil
}

func CrawlGoogle(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
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

	// Log do tamanho do HTML extraído
	log.Println("Tamanho do HTML extraído:", len(htmlContent))

	// Processar o HTML para extrair os produtos
	produtos, err := processarHTML(htmlContent)
	if err != nil {
		log.Println("Falha ao processar HTML:", err)
		return nil, fmt.Errorf("falha ao processar HTML: %v", err)
	}

	// Retornar os produtos extraídos
	return produtos, nil
}

func CrawlGoogle1111(query string) ([]Produto, error) {
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

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
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

	// Log do tamanho do HTML extraído
	log.Println("Tamanho do HTML extraído:", len(htmlContent))

	// Analisar o HTML e extrair os dados dos produtos
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao analisar o HTML:", err)
		return nil, fmt.Errorf("falha ao analisar o HTML: %v", err)
	}

	var produtos []Produto

	doc.Find(".pla-unit").Each(func(i int, s *goquery.Selection) {
		// Extrair dados do produto
		priceText := s.Find(".e10twf").Text()
		price, err := strconv.ParseFloat(strings.ReplaceAll(priceText, "$", ""), 64)
		if err != nil {
			price = 0.0 // Valor padrão se não conseguir converter
		}

		description := s.Find(".plantl.pla-unit-title-link").Text()
		source := s.Find(".zPEcBd").Text()
		link, _ := s.Find("a").Attr("href")
		imageURL, _ := s.Find("img").Attr("src")

		promoted := false // Lógica para determinar se o produto está em promoção pode ser adicionada aqui

		produto := Produto{
			Price:       price,
			Promotion:   promoted,
			Description: description,
			Source:      source,
			Link:        link,
			ImageURL:    imageURL,
		}

		produtos = append(produtos, produto)
	})

	// Retornar a lista de produtos extraídos
	return produtos, nil
}

func CrawlGoogleFuncionando(query string) ([]Produto, error) {
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

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
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

	// Log do tamanho do HTML extraído
	log.Println("Tamanho do HTML extraído:", len(htmlContent))

	// Log do HTML extraído para análise
	log.Println("HTML da página extraído:")
	log.Println(htmlContent)

	// Retornar uma lista vazia de produtos
	return []Produto{}, nil
}

func CrawlGoogleXXX(query string) ([]Produto, error) {
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

	var produtos []Produto

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(ctx, chromedp.Navigate(startURL))
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	for {
		// Esperar o carregamento da página com timeout específico
		log.Println("Esperando os resultados da página da coleta...")
		err = chromedp.Run(ctx, chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery))
		log.Println("Esperando os resultados da página da coleta depois de chromedp.Run ...")
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade dos resultados:", err)

			// Adicionando log para verificar o estado da página
			log.Println("Verificando o HTML da página após falha...")
			var htmlContent string
			err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
			if err != nil {
				log.Println("Erro ao extrair HTML:", err)
			} else {
				log.Println("HTML da página:", htmlContent) // Adicionando log para diagnóstico
			}

			break
		}
		log.Println("Visibilidade confirmada. Continuando...")

		// Extrair o HTML da página
		log.Println("Extraindo HTML da página...")
		var htmlContent string
		err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
		if err != nil {
			log.Println("Falha ao extrair HTML:", err)
			return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
		}

		// Log do tamanho do HTML extraído para verificar se está completo
		log.Println("Tamanho do HTML extraído:", len(htmlContent), "bytes")

		log.Println("HTML extraído com sucesso. Processando o HTML...")

		// Parsear o HTML com goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			log.Println("Falha ao parsear HTML:", err)
			return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
		}

		// Extrair detalhes dos produtos
		log.Println("Extraindo produtos da página...")
		doc.Find("div.sh-dgr__grid-result").Each(func(index int, item *goquery.Selection) {
			description := item.Find(".tAxDx").Text()
			priceText := item.Find(".a8Pemb").Text()
			log.Println("Raw price text:", priceText)

			price, err := formatarPreco(priceText)
			if err != nil {
				log.Println("Erro ao formatar o preço:", err)
				price = 0.0
			}
			log.Println("Formatted price:", price)

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
			log.Println("Produto encontrado:", produto)
		})

		// Verificar se há uma próxima página com timeout específico
		log.Println("Verificando se há uma próxima página...")
		var nextPageExists bool

		// Verificando se o botão de próxima página existe
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('a#pnnext') !== null`, &nextPageExists))
		if err != nil {
			log.Println("Erro ao verificar próxima página:", err)
			return nil, fmt.Errorf("erro ao verificar próxima página: %v", err)
		}

		if !nextPageExists {
			log.Println("Não há mais páginas para navegar.")
			break
		}

		log.Println("Próxima página encontrada, tentando navegar...")

		// Verificando se o elemento está visível
		err = chromedp.Run(ctx, chromedp.WaitVisible(`a#pnnext`, chromedp.ByQuery))
		if err != nil {
			log.Println("Erro ao esperar pela visibilidade do botão de próxima página:", err)
			return nil, fmt.Errorf("erro ao esperar pela visibilidade do botão de próxima página: %v", err)
		}

		// Navegar para a próxima página
		err = chromedp.Run(ctx, chromedp.Click(`a#pnnext`, chromedp.ByQuery, chromedp.NodeVisible))
		if err != nil {
			log.Println("Falha ao navegar para a próxima página:", err)
			return nil, fmt.Errorf("falha ao navegar para a próxima página: %v", err)
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		log.Println("Aguardando para evitar rate limiting...")
		time.Sleep(3 * time.Second)
	}

	// Log dos produtos coletados
	log.Println("Total de produtos coletados:", len(produtos))
	for _, produto := range produtos {
		log.Println("Produto:", produto)
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

func processarHTML(htmlContent string) ([]Produto, error) {
	// Criar um documento goquery a partir do HTML extraído
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao criar documento goquery:", err)
		return nil, fmt.Errorf("falha ao criar documento goquery: %v", err)
	}

	var produtos []Produto

	// Selecionar os contêineres de produtos (isso depende da estrutura HTML da página)
	doc.Find(".sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		// Extrair informações relevantes (ajustar seletores conforme a página)
		precoStr := s.Find(".a8Pemb").Text()
		descricao := s.Find(".shntl").Text()
		link, _ := s.Find("a").Attr("href")
		imagem, _ := s.Find("img").Attr("src")
		source := s.Find(".aULzUe").Text() // Exemplo de extração da loja/fonte

		// Conversão de preço para float64
		preco, err := strconv.ParseFloat(strings.ReplaceAll(precoStr, ",", ""), 64)
		if err != nil {
			log.Println("Erro ao converter preço:", err)
			preco = 0.0
		}

		// Criar e adicionar o produto à lista
		produto := Produto{
			Price:       preco,
			Description: descricao,
			Source:      source,
			Link:        "https://www.google.com" + link, // Adicionar domínio completo
			ImageURL:    imagem,
		}

		produtos = append(produtos, produto)
	})

	return produtos, nil
}
