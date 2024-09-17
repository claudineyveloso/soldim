package crawler

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
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

func CrawlGoogle(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3`),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

	// Criar um novo contexto do Chromium com as opções
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Definir o timeout
	ctx, cancel = context.WithTimeout(ctx, 3*time.Minute) // Ajuste o tempo conforme necessário
	defer cancel()

	// Criar o contexto para navegação
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery), // Aguarde o carregamento do corpo da página
		chromedp.Sleep(5*time.Second),                  // Ajuste o tempo de espera, se necessário
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	selecao := doc.Find(".tAxDx")
	if selecao.Length() == 0 {
		log.Println("Nenhum elemento com a classe tAxDx foi encontrado.")
	} else {
		selecao.Each(func(i int, s *goquery.Selection) {
			texto := s.Text() // Extrair o texto do elemento h3
			log.Printf("Elemento %d com classe tAxDx: %s\n", i, texto)
		})
	}

	var produtos []Produto
	doc.Find(".tAxDx").Each(func(i int, s *goquery.Selection) {
		htmlContent, err := s.Html()
		if err != nil {
			log.Printf("Erro ao obter HTML do elemento %d: %v\n", i, err)
			return
		}
		log.Printf("Elemento %d com classe tAxDx: %s\n", i, htmlContent)
	})

	log.Println("Produtos coletados:", len(produtos))
	// log.Println("Conteudo coletados:", htmlContent)

	return produtos, nil
}

func CrawlGoogleAVABV(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3`),
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
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.Sleep(5*time.Second), // Aumentar o tempo de espera
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
		chromedp.WaitVisible(".tAxDx", chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	// Log do tamanho do HTML extraído
	// log.Println("Tamanho do HTML extraído:", len(htmlContent))
	// log.Println("HTML extraído:\n", htmlContent) // Logar o HTML para verificar

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Verificar se a div esperada está presente
	if doc.Find(".sh-dgr__grid-result").Length() == 0 {
		log.Println("Não foram encontradas div.sh-dgr__grid-result")
	}

	var produtos []Produto
	// var mu sync.Mutex
	var wg sync.WaitGroup

	// Encontrar todas as divs de resultados
	// doc.Find("div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
	// 	wg.Add(1)
	// 	go func(s *goquery.Selection) {
	// 		defer wg.Done()
	// 		log.Println("Claudiney Veloso")
	// 		// Coletar dados do produto (a ser preenchido)
	// 	}(s)
	// 	log.Println("Total de Claudiney Veloso Coletado", doc.Find("div.sh-dgr__grid-result").Length())
	// })
	//
	// doc.Find("div").Each(func(i int, s *goquery.Selection) {
	// 	htmlContent, err := s.Html()
	// 	if err != nil {
	// 		log.Printf("Erro ao obter HTML do elemento %d: %v\n", i, err)
	// 		return
	// 	}
	// 	log.Printf("Elemento %d: %s\n", i, htmlContent)
	// })
	//
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
	// Esperar todas as goroutines terminarem
	wg.Wait()

	log.Println("Produtos coletados:", len(produtos))

	// Exibir os produtos coletados
	for _, produto := range produtos {
		log.Printf("Nome: %s\nLink: %s\nImagem: %s\nPreço: %.2f\n\n",
			produto.Description, produto.Link, produto.ImageURL, produto.Price)
	}

	// Retornar a lista de produtos coletados
	return produtos, nil
}

func CrawlGoogle888(query string) ([]Produto, error) {
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
	log.Println(htmlContent)

	err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	log.Println(doc.Find("div.sh-dgr__grid-result").Length())

	log.Println("Conteúdo HTML coletado")
	log.Println("Tamanho do HTML extraído:", len(htmlContent))
	// Retornar uma lista vazia de produtos
	return []Produto{}, nil
}

func CrawlGoogleCCC(query string) ([]Produto, error) {
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

	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Println("Iniciando visita:", startURL)

	// Navegar até a URL inicial
	err := chromedp.Run(timeoutCtx, chromedp.Navigate(startURL), chromedp.Sleep(5*time.Second), chromedp.WaitVisible(`div.sh-dgr__grid-result`, chromedp.ByQuery))
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

	// Coletar dados de até 10 produtos
	var produtos []Produto
	log.Println("Coletando produtos...")
	doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= 10 {
			return false // Interrompe a iteração após 10 produtos
		}
		log.Printf("Produto %d encontrado", i+1)

		nome := s.Find("h3.tAxDx").Text()
		link, _ := s.Find("a.xCpuod").Attr("href")
		imagemURL, _ := s.Find("img").Attr("src")
		precoStr := s.Find("span.a8Pemb").Text()

		// Capturar o texto visível do fornecedor usando JavaScript
		var fornecedor string
		err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector("div.aULzUe").innerText`, &fornecedor))
		if err != nil {
			log.Println("Falha ao extrair fornecedor:", err)
			fornecedor = "Desconhecido"
		}

		// Verificar se o preço foi encontrado
		if precoStr == "" {
			log.Println("Preço não encontrado para o produto:", nome)
			precoStr = "0" // Valor padrão caso o preço não seja encontrado
		}

		// Formatar o preço
		precoStr = strings.TrimSpace(precoStr)                // Remove espaços extras
		precoStr = strings.ReplaceAll(precoStr, "R$", "")     // Remove símbolo de moeda
		precoStr = strings.ReplaceAll(precoStr, ".", "")      // Remove pontos (milhares)
		precoStr = strings.ReplaceAll(precoStr, ",", ".")     // Substitui vírgulas por pontos (decimais)
		precoStr = strings.ReplaceAll(precoStr, "\u00a0", "") // Remove espaços não separáveis

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

	// Criar e salvar o arquivo CSV
	var filePath string
	if os.Getenv("ENV") == "production" { // Heroku
		filePath = "/tmp/produtos.csv"
	} else { // Ambiente local
		filePath = "produtos.csv"
	}
	log.Printf("Tentando criar o arquivo CSV em %s", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Falha ao criar o arquivo CSV:", err)
		return nil, fmt.Errorf("falha ao criar o arquivo CSV: %v", err)
	}
	defer file.Close()
	log.Printf("Arquivo CSV criado com sucesso em %s", filePath)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escrever cabeçalho
	err = writer.Write([]string{"Description", "Link", "ImageURL", "Price", "Source"})
	if err != nil {
		log.Println("Falha ao escrever cabeçalho no CSV:", err)
		return nil, fmt.Errorf("falha ao escrever cabeçalho no CSV: %v", err)
	}

	// Escrever os dados dos produtos
	for _, produto := range produtos {
		record := []string{
			produto.Description,
			produto.Link,
			produto.ImageURL,
			strconv.FormatFloat(produto.Price, 'f', 2, 64),
			produto.Source,
		}
		err = writer.Write(record)
		if err != nil {
			log.Println("Falha ao escrever registro no CSV:", err)
			return nil, fmt.Errorf("falha ao escrever registro no CSV: %v", err)
		}
	}

	log.Println("Arquivo CSV criado com sucesso")
	return produtos, nil
}

func CrawlGoogleFFF(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

	// Timeout ajustado para 30 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Criar contexto do Chromium com as opções
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
		chromedp.Sleep(8*time.Second), // Pode ser ajustado se necessário
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

	// Coletar dados de até 10 produtos
	var produtos []Produto
	doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= 5 {
			return false // Interrompe a iteração após 10 produtos
		}

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

	filePath := "produtos.csv"
	err = os.WriteFile(filePath, []byte(htmlContent), 0644)
	if err != nil {
		log.Println("Falha ao gravar HTML em arquivo:", err)
	} else {
		log.Println("HTML salvo com sucesso em", filePath)
	}

	// Retornar a lista de produtos coletados
	return produtos, nil
}
