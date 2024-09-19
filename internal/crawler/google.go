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
	"github.com/chromedp/cdproto/network"
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
	var execPath string
	if os.Getenv("HEROKU") == "true" {
		execPath = "/app/.chrome-for-testing/chrome-linux64/chrome"
	} else {
		execPath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}

	// Configurar opções para o Chromium
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(execPath),
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

	// Definir User-Agent para imitar o comportamento do Chrome em Linux
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.6668.58 Safari/537.36"

	// Variável para armazenar parte do HTML da página
	var htmlContent string
	err := chromedp.Run(ctx,
		network.SetExtraHTTPHeaders(network.Headers{
			"Accept-Language": "en-US,en;q=0.9",
		}),
		emulation.SetUserAgentOverride(userAgent),
		chromedp.Navigate(startURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
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

	log.Println("HTML recebido:", htmlContent[:40000])

	// var produtos []Produto
	//
	// // Usar seletores diferentes para local e Heroku
	// descriptions := doc.Find("span.pymv4e, h3.tAxDx")
	// prices := doc.Find("span.lmQWe, span.a8Pemb")
	// prices1 := doc.Find("span.DoCHT")
	// sources := doc.Find("span.zPEcBd, .aULzUe.IuHnof")
	// images := doc.Find(".D6nsM, .ArOc1c img")
	//
	// // Adicionar logs para verificar o número de elementos encontrados
	// log.Printf("Número de descrições encontradas: %d", descriptions.Length())
	// log.Printf("Número de preços encontrados: %d", prices.Length())
	// log.Printf("Número de preços encontrados: %d", prices1.Length())
	// log.Printf("Número de fontes encontradas: %d", sources.Length())
	// log.Printf("Número de imagens encontradas: %d", images.Length())
	//
	// // Adicionar logs para verificar conteúdo de preços
	// // for i := 0; i < prices.Length(); i++ {
	// // 	price := prices.Eq(i).Text()
	// // 	log.Printf("Preço encontrado: %s", price)
	// // }
	//
	// // Verificar se a quantidade de produtos está de acordo com a quantidade de elementos descritos
	// productCount := descriptions.Length()
	// if productCount != prices.Length() || productCount != sources.Length() || productCount != images.Length() {
	// 	log.Println("Os números de descrições, preços, fontes e imagens não coincidem!")
	// 	return nil, fmt.Errorf("número inconsistente de produtos extraídos")
	// }
	//
	// // Iterar sobre os produtos encontrados e associar as informações
	// for i := 0; i < productCount; i++ {
	// 	produto := Produto{}
	//
	// 	// Coletar descrição
	// 	produto.Description = descriptions.Eq(i).Text()
	//
	// 	// Coletar e converter preço
	// 	priceStr := prices.Eq(i).Text()
	// 	priceStr = strings.ReplaceAll(priceStr, "R$", "")     // Remover o símbolo "R$"
	// 	priceStr = strings.ReplaceAll(priceStr, "\u00a0", "") // Remover o espaço não separável
	// 	priceStr = strings.ReplaceAll(priceStr, ",", ".")     // Substituir vírgulas por pontos
	// 	priceStr = strings.TrimSpace(priceStr)                // Remover espaços em branco extras
	//
	// 	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	// 	if err != nil {
	// 		log.Println("Erro ao converter o preço:", err)
	// 		priceFloat = 0.0 // Definir um valor padrão em caso de erro
	// 	}
	// 	produto.Price = priceFloat
	//
	// 	// Coletar fonte (Source) e extrair texto após a última chave }
	// 	sourceText := sources.Eq(i).Text()
	//
	// 	// Adicionar a lógica de extração
	// 	lastBraceIndex := strings.LastIndex(sourceText, "}")
	// 	if lastBraceIndex != -1 {
	// 		textAfterCSS := sourceText[lastBraceIndex+1:]
	// 		textAfterCSS = strings.TrimSpace(textAfterCSS) // Remover espaços em branco
	// 		produto.Source = textAfterCSS
	// 	} else {
	// 		produto.Source = sourceText
	// 	}
	//
	// 	// Coletar imagem
	// 	if imgSrc, exists := images.Eq(i).Attr("src"); exists {
	// 		produto.ImageURL = imgSrc
	// 	} else if imgDataSrc, exists := images.Eq(i).Attr("data-src"); exists {
	// 		produto.ImageURL = imgDataSrc
	// 	}
	//
	// 	if doc.Find(".sh-dgr__grid-result").Length() > 0 {
	// 		log.Println("A classe .sh-dgr__grid-result foi encontrada.")
	// 	} else {
	// 		log.Println("A classe .sh-dgr__grid-result não foi encontrada.")
	// 	}
	//
	// 	// produtos = append(produtos, produto)
	// }
	// log.Println("Conteúdo HTML coletado:", htmlContent)

	// Retornar a lista de produtos
	return nil, nil
	// return produtos, nil
}
