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
		nome := s.Find("h3.tAxDx").Text()
		link, _ := s.Find("a.xCpuod").Attr("href")
		imagemURL, _ := s.Find("img").Attr("src")
		precoTexto := s.Find("span.a8Pemb").Text()
		fornecedor := s.Find("div.aULzUe").Text()

		// Limpar o texto do preço
		precoTexto = strings.ReplaceAll(precoTexto, "R$", "") // Remover o símbolo de moeda
		precoTexto = strings.ReplaceAll(precoTexto, ".", "")  // Remover separadores de milhar
		precoTexto = strings.ReplaceAll(precoTexto, ",", ".") // Substituir vírgula decimal por ponto

		// Converter o preço para float64
		preco, err := strconv.ParseFloat(strings.TrimSpace(precoTexto), 64)
		if err != nil {
			log.Printf("Falha ao converter o preço %s: %v", precoTexto, err)
			preco = 0.0 // Definir preço como 0.0 em caso de erro
		}

		// Montar estrutura Produto com os dados coletados
		produto := Produto{
			Description: nome,
			Link:        "https://www.google.com" + link,
			ImageURL:    imagemURL,
			Price:       preco, // Preço agora é um float64
			Source:      fornecedor,
		}

		produtos = append(produtos, produto)
	})

	// Verificar se produtos foram coletados corretamente
	if len(produtos) == 0 {
		log.Println("Nenhum produto encontrado.")
		return nil, fmt.Errorf("nenhum produto encontrado")
	}

	// Retornar a lista de produtos
	return produtos, nil
}
