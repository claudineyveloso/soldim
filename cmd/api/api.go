package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/services/deposit"
	depositbling "github.com/claudineyveloso/soldim.git/internal/services/deposit_bling"
	depositproduct "github.com/claudineyveloso/soldim.git/internal/services/deposit_product"
	"github.com/claudineyveloso/soldim.git/internal/services/draft"
	generatetoken "github.com/claudineyveloso/soldim.git/internal/services/generate_token"
	"github.com/claudineyveloso/soldim.git/internal/services/healthy"
	"github.com/claudineyveloso/soldim.git/internal/services/product"
	productbling "github.com/claudineyveloso/soldim.git/internal/services/product_bling"
	refreshtoken "github.com/claudineyveloso/soldim.git/internal/services/refresh_token"
	saleschannel "github.com/claudineyveloso/soldim.git/internal/services/sales_channel"
	saleschannelbling "github.com/claudineyveloso/soldim.git/internal/services/sales_channel_bling"
	"github.com/claudineyveloso/soldim.git/internal/services/search"
	searchresult "github.com/claudineyveloso/soldim.git/internal/services/search_result"
	"github.com/claudineyveloso/soldim.git/internal/services/stock"
	supplierproduct "github.com/claudineyveloso/soldim.git/internal/services/supplier_product"
	"github.com/claudineyveloso/soldim.git/internal/services/token"
	"github.com/claudineyveloso/soldim.git/internal/services/user"
	"github.com/claudineyveloso/soldim.git/internal/services/webhook"
	"github.com/gorilla/handlers"
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
	healthy.RegisterRoutes(r)
	webhook.RegisterRoutes(r)
	generatetoken.RegisterRoutes(r)
	refreshtoken.RegisterRoutes(r)
	productbling.RegisterRoutes(r)
	saleschannelbling.RegisterRoutes(r)
	depositbling.RegisterRoutes(r)
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)

	searchStore := search.NewStore(s.db)
	searchresultStore := searchresult.NewStore(s.db)
	draftStore := draft.NewStore(s.db)
	searchHandler := search.NewHandler(searchStore, searchresultStore, draftStore)
	searchHandler.RegisterRoutes(r)

	draftHandler := draft.NewHandler(draftStore)
	draftHandler.RegisterRoutes(r)

	// searchresultStore := searchresult.NewStore(s.db)
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

	fmt.Println("Server started on http://localhost:8080")
	// return http.ListenAndServe("localhost:8080", r)
	return http.ListenAndServe("localhost:8080",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(r))
}
