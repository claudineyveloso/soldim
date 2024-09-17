package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
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
		chromedp.Sleep(10*time.Second), // Aumentar o tempo de espera
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	// Log do tamanho do HTML extraído
	log.Println("Tamanho do HTML extraído:", len(htmlContent))
	log.Println("HTML extraído:\n", htmlContent) // Logar o HTML para verificar

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	var produtos []Produto
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Encontrar todos os resultados no Heroku e no ambiente local
	doc.Find("div.gkQHve, div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			defer wg.Done()

			// Coletar dados do produto, ajustando para os diferentes seletores
			nome := s.Find(".gkQHve, .tAxDx").Text()
			link, _ := s.Find("a.xCpuod").Attr("href")
			imagemURL, _ := s.Find("img.VeBrne, img").Attr("src")
			precoStr := s.Find("span.lmQWe, span.a8Pemb").Text()
			source := s.Find("span.WJMUdc").Text()

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
				preco = 0
			}

			produto := Produto{
				Description: nome,
				Link:        "https://www.google.com" + link,
				ImageURL:    imagemURL,
				Price:       preco,
				Source:      source,
			}

			mu.Lock()
			produtos = append(produtos, produto)
			mu.Unlock()
		}(s)
	})

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

func CrawlGoogleRRR(query string) ([]Produto, error) {
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
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.Sleep(10*time.Second), // Aumentar o tempo de espera
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao iniciar a visita:", err)
		return nil, fmt.Errorf("falha ao iniciar a visita: %v", err)
	}

	// Log do tamanho do HTML extraído
	log.Println("Tamanho do HTML extraído:", len(htmlContent))
	log.Println("HTML extraído:\n", htmlContent) // Logar o HTML para verificar

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Verificar se a div esperada está presente
	if doc.Find("div.sh-dgr__grid-result").Length() == 0 {
		log.Println("Não foram encontradas div.sh-dgr__grid-result")
	}

	var produtos []Produto
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Encontrar todas as divs de resultados
	doc.Find("div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			defer wg.Done()

			// Coletar dados do produto
			nome := s.Find("h3.tAxDx").Text()
			link, _ := s.Find("a.xCpuod").Attr("href")
			imagemURL, _ := s.Find("img").Attr("src")
			precoStr := s.Find("span.a8Pemb").Text()

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
				preco = 0
			}

			produto := Produto{
				Description: nome,
				Link:        "https://www.google.com" + link,
				ImageURL:    imagemURL,
				Price:       preco,
			}

			mu.Lock()
			produtos = append(produtos, produto)
			mu.Unlock()
		}(s)
	})

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

func CrawlGoogleTTT(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36`),
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

	// selecao := doc.Find(".tAxDx")
	// if selecao.Length() == 0 {
	// 	log.Println("Nenhum elemento com a classe tAxDx foi encontrado.")
	// } else {
	// 	selecao.Each(func(i int, s *goquery.Selection) {
	// 		texto := s.Text() // Extrair o texto do elemento h3
	// 		log.Printf("Elemento %d com classe tAxDx: %s\n", i, texto)
	// 	})
	// }
	//
	// selecaoProd := doc.Find(".gkQHve")
	// if selecaoProd.Length() == 0 {
	// 	log.Println("Nenhum elemento com a classe gkQHve foi encontrado.")
	// } else {
	// 	selecaoProd.Each(func(i int, s *goquery.Selection) {
	// 		texto := s.Text() // Extrair o texto do elemento h3
	// 		log.Printf("Elemento %d com classe gkQHve: %s\n", i, texto)
	// 	})
	// }
	//
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
	log.Println("Conteudo coletados:", htmlContent)

	return produtos, nil
}
