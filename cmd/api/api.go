package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/claudineyveloso/soldim.git/internal/services/contact"
	contactbling "github.com/claudineyveloso/soldim.git/internal/services/contact_bling"
	"github.com/claudineyveloso/soldim.git/internal/services/deposit"
	depositbling "github.com/claudineyveloso/soldim.git/internal/services/deposit_bling"
	depositproduct "github.com/claudineyveloso/soldim.git/internal/services/deposit_product"
	"github.com/claudineyveloso/soldim.git/internal/services/draft"
	generatetoken "github.com/claudineyveloso/soldim.git/internal/services/generate_token"
	"github.com/claudineyveloso/soldim.git/internal/services/healthy"
	"github.com/claudineyveloso/soldim.git/internal/services/heroku"
	itemssalesorder "github.com/claudineyveloso/soldim.git/internal/services/items_sales_order"
	"github.com/claudineyveloso/soldim.git/internal/services/product"
	productbling "github.com/claudineyveloso/soldim.git/internal/services/product_bling"
	productssalesorder "github.com/claudineyveloso/soldim.git/internal/services/products_sales_order"
	refreshtoken "github.com/claudineyveloso/soldim.git/internal/services/refresh_token"
	saleschannel "github.com/claudineyveloso/soldim.git/internal/services/sales_channel"
	saleschannelbling "github.com/claudineyveloso/soldim.git/internal/services/sales_channel_bling"
	salesorder "github.com/claudineyveloso/soldim.git/internal/services/sales_order"
	salesorderbling "github.com/claudineyveloso/soldim.git/internal/services/sales_order_bling"
	"github.com/claudineyveloso/soldim.git/internal/services/search"
	searchresult "github.com/claudineyveloso/soldim.git/internal/services/search_result"
	"github.com/claudineyveloso/soldim.git/internal/services/situation"
	"github.com/claudineyveloso/soldim.git/internal/services/stock"
	"github.com/claudineyveloso/soldim.git/internal/services/store"
	supplierproduct "github.com/claudineyveloso/soldim.git/internal/services/supplier_product"
	"github.com/claudineyveloso/soldim.git/internal/services/token"
	"github.com/claudineyveloso/soldim.git/internal/services/triage"
	"github.com/claudineyveloso/soldim.git/internal/services/user"
	webhooksales "github.com/claudineyveloso/soldim.git/internal/services/webhook_sales"
	webhookstock "github.com/claudineyveloso/soldim.git/internal/services/webhook_stock"
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	healthy.RegisterRoutes(r)
	heroku.RegisterRoutes(r)
	generatetoken.RegisterRoutes(r)
	refreshtoken.RegisterRoutes(r)
	productbling.RegisterRoutes(r)
	saleschannelbling.RegisterRoutes(r)
	depositbling.RegisterRoutes(r)
	contactbling.RegisterRoutes(r)
	salesorderbling.RegisterRoutes(r)
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)

	webhookstock.RegisterRoutes(r)
	webhooksales.RegisterRoutes(r)

	searchStore := search.NewStore(s.db)
	searchresultStore := searchresult.NewStore(s.db)
	draftStore := draft.NewStore(s.db)
	searchHandler := search.NewHandler(searchStore, searchresultStore, draftStore)
	searchHandler.RegisterRoutes(r)

	draftHandler := draft.NewHandler(draftStore)
	draftHandler.RegisterRoutes(r)

	searchresultHandler := searchresult.NewHandler(searchresultStore)
	searchresultHandler.RegisterRoutes(r)

	tokenStore := token.NewStore(s.db)
	tokenHandler := token.NewHandler(tokenStore)
	tokenHandler.RegisterRoutes(r)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(r)

	salesChannelStore := saleschannel.NewStore(s.db)
	salesChannelHandler := saleschannel.NewHandler(salesChannelStore)
	salesChannelHandler.RegisterRoutes(r)

	depositStore := deposit.NewStore(s.db)
	depositHandler := deposit.NewHandler(depositStore)
	depositHandler.RegisterRoutes(r)

	stockStore := stock.NewStore(s.db)
	stockHandler := stock.NewHandler(stockStore)
	stockHandler.RegisterRoutes(r)

	depositproductStore := depositproduct.NewStore(s.db)
	depositproductHandler := depositproduct.NewHandler(depositproductStore)
	depositproductHandler.RegisterRoutes(r)

	supplierproductStore := supplierproduct.NewStore(s.db)
	supplierproductHandler := supplierproduct.NewHandler(supplierproductStore)
	supplierproductHandler.RegisterRoutes(r)

	situationStore := situation.NewStore(s.db)
	situationHandler := situation.NewHandler(situationStore)
	situationHandler.RegisterRoutes(r)

	storeStore := store.NewStore(s.db)
	storeHandler := store.NewHandler(storeStore)
	storeHandler.RegisterRoutes(r)

	salesOrderStore := salesorder.NewStore(s.db)
	salesOrderHandler := salesorder.NewHandler(salesOrderStore)
	salesOrderHandler.RegisterRoutes(r)

	itemsSalesOrderStore := itemssalesorder.NewStore(s.db)
	itemsSalesOrderHandler := itemssalesorder.NewHandler(itemsSalesOrderStore)
	itemsSalesOrderHandler.RegisterRoutes(r)

	productSalesOrderStore := productssalesorder.NewStore(s.db)
	productSalesOrderHandler := productssalesorder.NewHandler(productSalesOrderStore)
	productSalesOrderHandler.RegisterRoutes(r)

	triageStore := triage.NewStore(s.db)
	triateHandler := triage.NewHandler(triageStore)
	triateHandler.RegisterRoutes(r)

	contactStore := contact.NewStore(s.db)
	contactHandler := contact.NewHandler(contactStore)
	contactHandler.RegisterRoutes(r)

	env := os.Getenv("ENVIRONMENT")
	var address string

	if env == "prod" {
		// Use a porta fornecida pelo Heroku
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
		address = ":" + port
	} else {
		// Se estiver em desenvolvimento, use o localhost
		address = "localhost:8080"
		fmt.Println("Server started on http://localhost:8080")
	}

	// Configurar CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Permite todas as origens, ajuste conforme necessário
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}).Handler(r)

	// Iniciar o servidor
	return http.ListenAndServe(address, corsHandler)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Loga o método e o caminho da requisição
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		// Passa a requisição para o próximo handler
		next.ServeHTTP(w, r)
	})
}
