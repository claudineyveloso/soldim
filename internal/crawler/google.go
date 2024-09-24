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
	"github.com/gocolly/colly"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
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
	const (
		chromeDriverPath = "/usr/bin/chromedriver"
		// chromeDriverPath = "bin/chromedriver" // Atualize com o caminho correto
		port = 8081
	)

	// Iniciar o serviço do ChromeDriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Erro ao iniciar o ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Definir capacidades do navegador
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"args": []string{
				// "--headless", // Executar em modo headless
				// "--no-sandbox",
				// "--disable-dev-shm-usage",
				"--headless",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--disable-gpu",
				"--enable-logging",
				"--v=1",
				"--log-level=ALL",              // Mais detalhes sobre o log
				"--remote-debugging-port=9222", // Para depuração
			},
		},
	}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Erro ao conectar ao WebDriver: %v", err)
	}
	defer wd.Quit()

	// Codificar a query string para ser usada na URL
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	log.Printf("Iniciando visita: %s", startURL)

	// Navegar até a URL inicial
	if err := wd.Get(startURL); err != nil {
		log.Printf("Falha ao iniciar a visita: %v", err)
		// return nil, log.Println("falha ao iniciar a visita: %v", err)
	}

	// Aguardar um tempo para evitar problemas com rate limiting
	time.Sleep(4000 * time.Millisecond)

	var produtos []Produto

	for {
		// Aguardar a visibilidade dos resultados
		var elementos []selenium.WebElement
		for start := time.Now(); time.Since(start) < 10*time.Second; {
			elementos, err = wd.FindElements(selenium.ByCSSSelector, "div.sh-dgr__grid-result")
			if err == nil && len(elementos) > 0 {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil || len(elementos) == 0 {
			log.Printf("Erro ao esperar pela visibilidade dos resultados: %v", err)
			break
		}

		// Extrair detalhes dos produtos
		for _, elemento := range elementos {
			var description, imageURL, rawURL, source string
			price := 0.0
			promotion := false

			// Extrair descrição
			descriptionElem, err := elemento.FindElement(selenium.ByCSSSelector, ".tAxDx")
			if err == nil {
				descText, _ := descriptionElem.Text()
				description = strings.TrimSpace(descText)
			}

			// Extrair preço
			priceElem, err := elemento.FindElement(selenium.ByCSSSelector, ".a8Pemb")
			if err == nil {
				priceText, _ := priceElem.Text()
				price = formatarPreco(priceText)
			}

			// Extrair link
			linkElem, err := elemento.FindElement(selenium.ByCSSSelector, "a")
			if err == nil {
				rawURL, _ = linkElem.GetAttribute("href") // Corrigido para GetAttribute
			}

			// Extrair imagem
			imageElem, err := elemento.FindElement(selenium.ByCSSSelector, ".ArOc1c img")
			if err == nil {
				imageURL, _ = imageElem.GetAttribute("src") // Corrigido para GetAttribute
			}

			// Extrair promoção
			promotionTextElem, err := elemento.FindElement(selenium.ByCSSSelector, ".fAcMNb span.Ib8pOd")
			if err == nil {
				promotionText, _ := promotionTextElem.Text()
				promotion = strings.TrimSpace(promotionText) == "PROMOÇÃO"
			}

			// Extrair fonte
			sourceElem, err := elemento.FindElement(selenium.ByCSSSelector, ".aULzUe")
			if err == nil {
				sourceText, _ := sourceElem.Text()
				source = strings.TrimSpace(sourceText)
			}

			// Montando o produto
			produto := Produto{
				Description: description,
				Price:       price,
				Source:      source,
				Link:        rawURL,
				ImageURL:    imageURL,
				Promotion:   promotion,
			}
			produtos = append(produtos, produto)
			log.Printf("Produto encontrado: %+v\n", produto)
		}

		// Verificar se há uma próxima página
		nextPageExists, err := wd.FindElement(selenium.ByCSSSelector, "a#pnnext")
		if err == nil {
			if err := nextPageExists.Click(); err != nil {
				log.Printf("Erro ao navegar para a próxima página: %v", err)
			}
		} else {
			break // Não há mais páginas
		}

		// Aguardar um tempo para evitar problemas com rate limiting
		time.Sleep(4000 * time.Millisecond)

		log.Printf("Total de produtos coletados: %d", len(produtos))
		for _, produto := range produtos {
			log.Printf("Produto: %+v", produto)
		}
	}

	return produtos, nil
}

func CrawlGoogleXX(query string) ([]Produto, error) {
	// Caminho para o chromedriver e porta
	const (
		chromeDriverPath = "bin/chromedriver" // Substitua pelo caminho para o seu chromedriver
		port             = 8081
	)

	// Iniciar o serviço do chromedriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Erro ao iniciar o serviço do ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Definir as capacidades do navegador
	caps := selenium.Capabilities{"browserName": "chrome"}

	// Configurar o Chrome para rodar em modo headless
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--window-size=1920,1080",
		},
	}
	caps.AddChrome(chromeCaps)

	// Conectar ao navegador
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Erro ao conectar ao WebDriver: %v", err)
	}
	defer wd.Quit()

	// Acessar a página do Google Shopping com a query desejada
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", url.QueryEscape(query))
	if err := wd.Get(startURL); err != nil {
		log.Fatalf("Erro ao acessar a página: %v", err)
	}

	time.Sleep(5 * time.Second)

	// productContainers, err := wd.FindElements(selenium.ByCSSSelector, "div.sh-dgr__grid-result")
	// if err != nil {
	// 	log.Printf("Erro ao buscar contêineres de produto: %v", err)
	// 	return nil, err
	// }
	//
	// // Aqui você pode continuar com a lógica de extração de dados dos produtos.
	// var produtos []Produto
	// for _, productElement := range productContainers {
	// 	produto := Produto{}
	//
	// 	// Coletar o nome do produto
	// 	nameElement, err := productElement.FindElement(selenium.ByCSSSelector, "h4.Xjkr3b")
	// 	if err == nil {
	// 		produto.Description, _ = nameElement.Text()
	// 	}
	//
	// 	// Coletar o preço do produto
	// 	priceElement, err := productElement.FindElement(selenium.ByCSSSelector, "span.a8Pemb")
	// 	if err == nil {
	// 		priceText, _ := priceElement.Text()
	// 		priceText = strings.ReplaceAll(priceText, "$", "")
	// 		priceText = strings.ReplaceAll(priceText, "R$", "")
	// 		priceText = strings.ReplaceAll(priceText, ".", "")
	// 		priceText = strings.ReplaceAll(priceText, ",", ".")
	// 		produto.Price, _ = strconv.ParseFloat(priceText, 64)
	// 	}
	//
	// 	// Coletar a fonte do produto (loja)
	// 	sourceElement, err := productElement.FindElement(selenium.ByCSSSelector, "span.aULzUe")
	// 	if err == nil {
	// 		produto.Source, _ = sourceElement.Text()
	// 	}
	//
	// 	// Coletar a URL da imagem do produto
	// 	imageElement, err := productElement.FindElement(selenium.ByCSSSelector, "img.sh-dgr__image")
	// 	if err == nil {
	// 		produto.ImageURL, _ = imageElement.GetAttribute("src")
	// 	}
	//
	// 	// Adicionar o produto ao slice
	// 	produtos = append(produtos, produto)
	//
	// 	// Logar o produto para verificação
	// 	log.Printf("Produto coletado: Nome: %s, Preço: %.2f, Fonte: %s, Imagem: %s",
	// 		produto.Description, produto.Price, produto.Source, produto.ImageURL)
	// }
	// Sua lógica de coleta de dados dos produtos...
	// Obter o HTML da página atual
	html, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Erro ao obter o HTML da página: %v", err)
	}

	// Salvar o HTML completo em um arquivo
	file, err := os.Create("pagina_completa_selenium.html")
	if err != nil {
		log.Printf("Erro ao criar arquivo: %v", err)
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		log.Printf("Erro ao escrever no arquivo: %v", err)
		return nil, err
	}

	log.Println("HTML salvo em pagina_completa_selenium.html")

	// Aqui você pode continuar com a lógica de extração de dados dos produtos.
	var produtos []Produto
	// Sua lógica de coleta de dados dos produtos...

	return produtos, nil
}

// func CrawlGoogleSelenium(query string) ([]Produto, error) {
func CrawlGoogleSelenium(query string) ([]Produto, error) {
	// Caminho para o chromedriver e porta
	const (
		seleniumPath     = "bin/selenium-server-standalone-3.5.0.jar"
		chromeDriverPath = "bin/chromedriver" // Substitua pelo caminho para o seu chromedriver
		port             = 8081
	)

	// Iniciar o serviço do ChromeDriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Conectar ao Chrome via WebDriver
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Erro ao conectar ao WebDriver: %v", err)
	}
	defer driver.Quit()

	// Acessar a página do Google Shopping com a query
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)
	if err := driver.Get(startURL); err != nil {
		log.Fatalf("Erro ao acessar a página: %v", err)
	}

	// Esperar alguns segundos para garantir que o conteúdo dinâmico foi carregado
	time.Sleep(5 * time.Second)

	// Coletar os elementos de produtos
	productElements, err := driver.FindElements(selenium.ByCSSSelector, "div.Ez5pwe")
	if err != nil {
		log.Fatalf("Erro ao buscar elementos de produtos: %v", err)
	}

	var produtos []Produto

	// Loop através dos elementos de produto e coletar as informações
	for _, productElement := range productElements {
		produto := Produto{}

		// Nome do produto
		descriptionElement, err := productElement.FindElement(selenium.ByCSSSelector, "div.gkQHve")
		if err == nil {
			produto.Description, _ = descriptionElement.Text()
		}

		// Preço do produto
		priceElement, err := productElement.FindElement(selenium.ByCSSSelector, "span.lmQWe")
		if err == nil {
			precoStr, _ := priceElement.Text()
			precoStr = strings.ReplaceAll(precoStr, "$", "")
			precoStr = strings.ReplaceAll(precoStr, "R$", "")
			precoStr = strings.ReplaceAll(precoStr, ".", "")
			precoStr = strings.ReplaceAll(precoStr, ",", ".")
			produto.Price, _ = strconv.ParseFloat(precoStr, 64)
		}

		// Fonte do produto (loja)
		sourceElement, err := productElement.FindElement(selenium.ByCSSSelector, "span.WJMUdc")
		if err == nil {
			produto.Source, _ = sourceElement.Text()
		}

		// URL da imagem
		imageElement, err := productElement.FindElement(selenium.ByCSSSelector, "img")
		if err == nil {
			produto.ImageURL, _ = imageElement.GetAttribute("src")
		}

		// Adicionar o produto à lista
		produtos = append(produtos, produto)

		// Logar o produto coletado
		log.Printf("Produto coletado: Nome: %s, Preço: %f, Fonte: %s, Imagem: %s", produto.Description, produto.Price, produto.Source, produto.ImageURL)
	}

	return produtos, nil
}

func CrawlGoogleColly(query string) ([]Produto, error) {
	encodedQuery := url.QueryEscape(query)
	startURL := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=shop", encodedQuery)

	// Criar uma nova instância do Colly com limitações
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
		colly.MaxDepth(2), // Limitar profundidade para evitar loops
	)

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
