package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
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

func CrawlGoogleUUU(query string) ([]Produto, error) {
	// Configurar contexto do chromedp com timeout de 10 minutos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções necessárias para Heroku
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Novo contexto para o navegador
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Variável para armazenar o HTML da página
	var htmlContent string

	// Codificar a query e construir a URL inicial
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", url.QueryEscape(query))
	log.Println("Iniciando visita:", startURL)

	// Navegar para a página de resultados de shopping
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.WaitVisible(`div.sh-dgr__grid-result`), // Esperar a página carregar
		chromedp.Sleep(5*time.Second),                   // Tempo extra para garantir carregamento
		chromedp.OuterHTML("body", &htmlContent),        // Coletar HTML da página
	)
	if err != nil {
		log.Println("Falha ao carregar a página de resultados:", err)
		return nil, fmt.Errorf("falha ao carregar a página de resultados: %v", err)
	}

	// Log do HTML coletado (parcial)
	log.Println("HTML coletado:", htmlContent[:500]) // Logar os primeiros 500 caracteres

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Inicializar slice de produtos
	var produtos []Produto

	// Coletar dados de cada produto
	doc.Find("div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		// Extrair a descrição do produto
		descricao := s.Find("div.EI11Pd h3.tAxDx").Text()
		log.Printf("Produto %d: %s\n", i+1, descricao) // Log da descrição do produto

		// Adicionar produto ao slice (por enquanto vazio)
		produtos = append(produtos, Produto{})
	})

	// Verificar se produtos foram coletados corretamente
	if len(produtos) == 0 {
		log.Println("Nenhum produto encontrado.")
		return produtos, nil // Retorne uma lista vazia ao invés de um erro
	}

	// Retornar a lista de produtos
	return produtos, nil
}

func CrawlGoogle(query string) ([]Produto, error) {
	// Configurar contexto do chromedp com timeout de 10 minutos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções necessárias para Heroku
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Novo contexto para o navegador
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Variável para armazenar o HTML da página
	var htmlContent string

	// Codificar a query e construir a URL inicial
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", url.QueryEscape(query))
	log.Println("Iniciando visita:", startURL)

	// Navegar para a página de resultados de shopping
	// err := chromedp.Run(ctx,
	// 	chromedp.Emulate(emulation.SetUserAgentOverride("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")),
	// 	chromedp.Navigate(startURL),
	// 	chromedp.WaitVisible(`div.sh-dgr__grid-result`), // Esperar a página carregar
	// 	chromedp.Sleep(5*time.Second),
	// 	chromedp.OuterHTML("body", &htmlContent), // Coletar HTML da página
	// )
	//

	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Navigate(startURL),
		chromedp.WaitVisible(`div.sh-dgr__grid-result`),
		chromedp.OuterHTML("body", &htmlContent),
	)
	if err != nil {
		log.Println("Falha ao carregar a página de resultados:", err)
		return nil, fmt.Errorf("falha ao carregar a página de resultados: %v", err)
	}

	// Log para verificar o HTML coletado
	log.Println("HTML coletado:", htmlContent[:2000]) // Exibe os primeiros 500 caracteres para verificar se está correto

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	err = os.WriteFile("/tmp/html_output.txt", []byte(htmlContent), 0644)
	if err != nil {
		log.Println("Falha ao gravar HTML em arquivo:", err)
	}

	// Inicializar slice de produtos (vazio neste caso)
	var produtos []Produto

	// Coletar dados de cada produto (somente log por enquanto)
	doc.Find("div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		log.Println("Produto encontrado.")
		// Lógica de extração de produtos será adicionada depois
	})

	// Log indicando que não estamos retornando produtos ainda
	log.Println("Retornando lista de produtos vazia.")

	// Retornar lista vazia (sem erro)
	return produtos, nil
}

func CrawlGoogleXXX(query string) ([]Produto, error) {
	// Configurar contexto do chromedp com timeout de 10 minutos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Criar um novo contexto do Chromium com as opções necessárias para Heroku
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)
	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Novo contexto para o navegador
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Variável para armazenar o HTML da página
	var htmlContent string

	// Codificar a query e construir a URL inicial
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", url.QueryEscape(query))
	log.Println("Iniciando visita:", startURL)

	// Navegar para a página de resultados de shopping
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.WaitVisible(`div.sh-dgr__grid-result`), // Esperar a página carregar
		chromedp.OuterHTML("body", &htmlContent),        // Coletar HTML da página
	)
	if err != nil {
		log.Println("Falha ao carregar a página de resultados:", err)
		return nil, fmt.Errorf("falha ao carregar a página de resultados: %v", err)
	}

	// Parsear o HTML usando goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Inicializar slice de produtos
	var produtos []Produto

	// Coletar dados de cada produto
	doc.Find("div.sh-dgr__grid-result").Each(func(i int, s *goquery.Selection) {
		log.Println("produtos encontrados.")
	})

	// Verificar se produtos foram coletados corretamente
	if len(produtos) == 0 {
		log.Println("Nenhum produto encontrado.")
		return nil, fmt.Errorf("nenhum produto encontrado")
	}

	// Retornar a lista de produtos
	return produtos, nil
}
