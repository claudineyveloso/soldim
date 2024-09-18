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
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

	// Definir contexto com timeout de 5 minutos para evitar longas execuções
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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

	// Variável para armazenar parte do HTML da página
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	if doc.Find(".sh-dgr__grid-result").Length() > 0 {
		log.Println("A classe .sh-dgr__grid-result foi encontrada.")
	} else {
		log.Println("A classe .sh-dgr__grid-result não foi encontrada.")
	}

	log.Println("###################################################################.")
	log.Println("Claudiney Veloso.")
	log.Println("###################################################################.")

	var description, price, source, image string

	doc.Find("span.pymv4e, h3.tAxDx").Each(func(i int, s *goquery.Selection) {
		description += s.Text() + " "
	})

	doc.Find("span.lmQWe, span.a8Pemb").Each(func(i int, s *goquery.Selection) {
		price += s.Text() + " "
	})

	// Extraindo fontes (local e Heroku)
	doc.Find("span.zPEcBd, .aULzUe.IuHnof").Each(func(i int, s *goquery.Selection) {
		source += s.Text() + " "
	})

	doc.Find(".D6nsM, .ArOc1c").Each(func(i int, s *goquery.Selection) {
		if imgSrc, exists := s.Find("img").Attr("src"); exists {
			image += imgSrc + " "
		} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
			image += imgDataSrc + " "
		}
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// Primeiro tenta pegar o data-src
		if dataSrc, exists := s.Attr("data-src"); exists {
			fmt.Println("Valor do data-src:", dataSrc)
		} else if src, exists := s.Attr("src"); exists {
			// Se data-src não existir, tenta pegar o src
			fmt.Println("Valor do src:", src)
		} else {
			fmt.Println("Nenhum dos atributos 'data-src' ou 'src' foi encontrado!")
		}
	})

	// Logar o nome extraído
	log.Println("Descrição extraída:", description)
	log.Println("Preço extraído:", price)
	log.Println("Fonte extraída:", source)
	log.Println("Imagens extraídas:", image)

	// var produtos []Produto

	// produto := Produto{
	// 	Description: description,
	// 	Link:        "https://www.google.com" + link,
	// 	ImageURL:    image,
	// 	Price:       price,
	// 	Source:      source,
	// }
	// //

	// Logar o conteúdo HTML
	// log.Println("Conteúdo HTML coletado:", htmlContent)

	// Processar o HTML ou retornar
	return nil, nil
}

func CrawlGoogleBABA(query string) ([]Produto, error) {
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
	// var className string

	// Executar a navegação e extrair o nome da classe da div com role="navigation"
	// err := chromedp.Run(ctx,
	// 	// Navegar para a URL
	// 	chromedp.Navigate(startURL),
	// 	// Esperar até que o elemento com role="navigation" seja visível
	// 	chromedp.WaitVisible(`[role="navigation"]`, chromedp.ByQuery),
	// 	// Extrair o atributo "class" do elemento
	// 	chromedp.AttributeValue(`[role="navigation"]`, "class", &className, nil),
	// )
	// if err != nil {
	// 	log.Println("Falha ao extrair classe da div com role='navigation':", err)
	// 	return nil, fmt.Errorf("falha ao extrair classe: %v", err)
	// }
	//
	// // Exibir o nome da classe extraído
	// log.Println("Nome da classe da div com role='navigation':", className)

	// Extrair o HTML da página
	var htmlContent string
	err := chromedp.Run(ctx,
		// Reutilizando o contexto para extrair o HTML da página
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	// Parsear o HTML com goquery
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	// if err != nil {
	// 	log.Println("Falha ao parsear HTML:", err)
	// 	return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	// }

	// Variáveis para armazenar dados
	// var description, price, source, href, image string
	// var description, price, source, image string
	// var description string
	//
	// // Extraindo descrições dos produtos
	// doc.Find("span.pymv4e, h3.tAxDx").Each(func(i int, s *goquery.Selection) {
	// 	description += s.Text() + " "
	// })
	//
	// Extraindo preços
	// doc.Find("span.lmQWe, span.a8Pemb").Each(func(i int, s *goquery.Selection) {
	// 	price += s.Text() + " "
	// })
	//
	// // Extraindo fontes (vendedor)
	// doc.Find("span.zPEcBd.LnPkof, .aULzUe.IuHnof").Each(func(i int, s *goquery.Selection) {
	// 	source += s.Text() + " "
	// })
	//
	// // Extraindo URLs de imagens
	// doc.Find(".D6nsM, .ArOc1c").Each(func(i int, s *goquery.Selection) {
	// 	if imgSrc, exists := s.Find("img").Attr("src"); exists {
	// 		image += imgSrc + " "
	// 	} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
	// 		image += imgDataSrc + " "
	// 	}
	// })

	// doc.Find(".plantl .pla-unit-title-link, .shntl .sh-np__click-target").Each(func(i int, s *goquery.Selection) {
	// 	// Variável para armazenar o link (href)
	// 	var href string
	//
	// 	// Tentar pegar o atributo href do elemento
	// 	if link, exists := s.Attr("href"); exists {
	// 		href = link
	// 	}
	//
	// 	// Verificar se há uma imagem no mesmo elemento e pegar o src ou data-src
	// 	if imgSrc, exists := s.Find("img").Attr("src"); exists {
	// 		image += imgSrc + " "
	// 	} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
	// 		image += imgDataSrc + " "
	// 	}
	//
	// 	// Logar o link e a imagem extraídos
	// })

	// Logar os dados extraídos
	// log.Println("Descrição extraída:", description)
	// log.Println("Preço extraído:", price)
	// log.Println("Fonte extraída:", source)
	// log.Println("Imagens extraídas:", image)
	// log.Println("Link extraído:", href)
	// log.Println("Imagem extraída:", image)

	log.Println("Conteúdo HTML coletado")

	log.Println("Coletando o html", htmlContent)

	// Retornar nil, já que não estamos processando os produtos por enquanto
	return nil, nil
}

func CrawlGoogleCaca(query string) ([]Produto, error) {
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

	// Extrair o HTML da página
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(startURL),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}

	// Logar o conteúdo HTML
	// log.Println("Conteúdo HTML coletado:", htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Falha ao parsear HTML:", err)
		return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	}

	// Localizar e atribuir texto das tags à variável nome
	var description, priceStr, source, image string

	doc.Find("span.pymv4e, h3.tAxDx").Each(func(i int, s *goquery.Selection) {
		description += s.Text() + " "
	})

	doc.Find("span.lmQWe, span.a8Pemb").Each(func(i int, s *goquery.Selection) {
		priceStr += s.Text() + " "
	})

	// Extraindo fontes (local e Heroku)
	doc.Find("span.zPEcBd, .aULzUe.IuHnof").Each(func(i int, s *goquery.Selection) {
		source += s.Text() + " "
	})

	doc.Find(".D6nsM, .ArOc1c").Each(func(i int, s *goquery.Selection) {
		if imgSrc, exists := s.Find("img").Attr("src"); exists {
			image += imgSrc + " "
		} else if imgDataSrc, exists := s.Find("img").Attr("data-src"); exists {
			image += imgDataSrc + " "
		}
	})

	// Logar o nome extraído
	log.Println("Descrição extraída:", description)
	log.Println("Preço extraído:", priceStr)
	log.Println("Fonte extraída:", source)
	log.Println("Imagens extraídas:", image)
	// Retornar nil, já que não estamos processando os produtos por enquanto
	var produtos []Produto
	// doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	nome := s.Find("h3.tAxDx").Text()
	// 	link, _ := s.Find("a.xCpuod").Attr("href")
	// 	imagemURL, _ := s.Find("img").Attr("src")
	// 	precoStr := s.Find("span.a8Pemb").Text()
	// 	fornecedor := "Desconhecido"
	//
	// 	// Verificar se o preço foi encontrado
	if priceStr == "" {
		priceStr = "0" // Valor padrão caso o preço não seja encontrado
	}
	//
	// 	// Formatar o preço
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, "R$", "")
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	priceStr = strings.ReplaceAll(priceStr, ".", "")
	priceStr = strings.ReplaceAll(priceStr, ",", ".")
	priceStr = strings.ReplaceAll(priceStr, "\u00a0", "")

	// Converter preço de string para float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("Erro ao converter preço '%s' para float64: %v", priceStr, err)
		price = 0 // Valor padrão caso a conversão falhe
	}
	//
	produto := Produto{
		Description: description,
		// Link:        "https://www.google.com" + link,
		ImageURL: image,
		Price:    price,
		Source:   source,
	}

	produtos = append(produtos, produto)

	return produtos, nil
}

func CrawlGoogleBBB(query string) ([]Produto, error) {
	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),            // Necessário para Heroku
		chromedp.Flag("disable-dev-shm-usage", true), // Pode ajudar a evitar problemas de memória
		chromedp.Flag("disable-gpu", true),           // O Heroku não precisa de GPU
	)

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

	// Extrair o HTML da página
	var htmlContent string
	err := chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery))
	if err != nil {
		log.Println("Falha ao extrair HTML:", err)
		return nil, fmt.Errorf("falha ao extrair HTML: %v", err)
	}
	log.Println("Conteúdo HTML coletado")

	log.Println("Coletando o html", htmlContent)

	// Parsear o HTML usando goquery
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	// if err != nil {
	// 	log.Println("Falha ao parsear HTML:", err)
	// 	return nil, fmt.Errorf("falha ao parsear HTML: %v", err)
	// }

	// selecao := doc.Find(".tAxDx")
	// if selecao.Length() == 0 {
	// 	log.Println("Nenhum elemento com a classe .tAxDx foi encontrado.")
	// } else {
	// 	selecao.Each(func(i int, s *goquery.Selection) {
	// 		htmlContent, err := s.Html()
	// 		if err != nil {
	// 			log.Printf("Erro ao obter HTML do elemento %d: %v\n", i, err)
	// 			return
	// 		}
	// 		log.Printf("Elemento %d com classe tAxDx: %s\n", i, htmlContent)
	// 	})
	// }
	//
	// // Coletar dados de até 10 produtos
	// var produtos []Produto
	// doc.Find("div.sh-dgr__grid-result").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	nome := s.Find("h3.tAxDx").Text()
	// 	link, _ := s.Find("a.xCpuod").Attr("href")
	// 	imagemURL, _ := s.Find("img").Attr("src")
	// 	precoStr := s.Find("span.a8Pemb").Text()
	// 	fornecedor := "Desconhecido"
	//
	// 	// Verificar se o preço foi encontrado
	// 	if precoStr == "" {
	// 		precoStr = "0" // Valor padrão caso o preço não seja encontrado
	// 	}
	//
	// 	// Formatar o preço
	// 	precoStr = strings.TrimSpace(precoStr)
	// 	precoStr = strings.ReplaceAll(precoStr, "R$", "")
	// 	precoStr = strings.ReplaceAll(precoStr, ".", "")
	// 	precoStr = strings.ReplaceAll(precoStr, ",", ".")
	// 	precoStr = strings.ReplaceAll(precoStr, "\u00a0", "")
	//
	// 	// Converter preço de string para float64
	// 	preco, err := strconv.ParseFloat(precoStr, 64)
	// 	if err != nil {
	// 		log.Printf("Erro ao converter preço '%s' para float64: %v", precoStr, err)
	// 		preco = 0 // Valor padrão caso a conversão falhe
	// 	}
	//
	// 	produto := Produto{
	// 		Description: nome,
	// 		Link:        "https://www.google.com" + link,
	// 		ImageURL:    imagemURL,
	// 		Price:       preco,
	// 		Source:      fornecedor,
	// 	}
	//
	// 	produtos = append(produtos, produto)
	//	return true
	//})
	//

	// log.Println("Produtos coletados:", len(produtos))

	// Retornar a lista de produtos coletados
	// return produtos, nil
	return nil, nil
}
