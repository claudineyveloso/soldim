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

	// Variável para armazenar o nome da classe da div com role="navigation"
	var className string

	// Executar a navegação e extrair o nome da classe da div com role="navigation"
	err := chromedp.Run(ctx,
		// Navegar para a URL
		chromedp.Navigate(startURL),
		// Esperar até que o elemento com role="navigation" seja visível
		chromedp.WaitVisible(`[role="navigation"]`, chromedp.ByQuery),
		// Extrair o atributo "class" do elemento
		chromedp.AttributeValue(`[role="navigation"]`, "class", &className, nil),
	)
	if err != nil {
		log.Println("Falha ao extrair classe da div com role='navigation':", err)
		return nil, fmt.Errorf("falha ao extrair classe: %v", err)
	}

	// Exibir o nome da classe extraído
	log.Println("Nome da classe da div com role='navigation':", className)

	// Extrair o HTML da página
	var htmlContent string
	err = chromedp.Run(ctx,
		// Reutilizando o contexto para extrair o HTML da página
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	// Parsear o HTML com goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Variáveis para armazenar dados
	var descriptions, prices, sources, hrefs, images []string

	// Extraindo descrições dos produtos
	doc.Find("span.pymv4e, h3.tAxDx").Each(func(i int, s *goquery.Selection) {
		descriptions = append(descriptions, s.Text())
	})

	// Extraindo preços
	doc.Find("span.lmQWe, span.a8Pemb").Each(func(i int, s *goquery.Selection) {
		prices = append(prices, s.Text())
	})

	// Extraindo fontes (vendedor)
	doc.Find("span.zPEcBd.LnPkof, .aULzUe.IuHnof").Each(func(i int, s *goquery.Selection) {
		sources = append(sources, s.Text())
	})

	// Extraindo URLs de imagens
	doc.Find(".D6nsM, .ArOc1c").Each(func(i int, s *goquery.Selection) {
		if imgSrc, exists := s.Find("img").Attr("src"); exists {
			images = append(images, imgSrc)
		} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
			images = append(images, imgDataSrc)
		}
	})

	// Extraindo URLs de produtos
	doc.Find(".plantl .pla-unit-title-link, .shntl .sh-np__click-target").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			hrefs = append(hrefs, href)
		}
	})

	// Logar os dados extraídos
	log.Println("Descrições extraídas:", descriptions)
	log.Println("Preços extraídos:", prices)
	log.Println("Fontes extraídas:", sources)
	log.Println("Imagens extraídas:", images)
	log.Println("Links extraídos:", hrefs)

	// Retornar nil, já que não estamos processando os produtos por enquanto
	return nil, nil
}
