package api

import (
	"database/sql"
	"fmt"
	"net/http"

	generatetoken "github.com/claudineyveloso/soldim.git/internal/services/generate_token"
	"github.com/claudineyveloso/soldim.git/internal/services/healthy"
	"github.com/claudineyveloso/soldim.git/internal/services/product"
	salechannel "github.com/claudineyveloso/soldim.git/internal/services/sale_channel"
	"github.com/claudineyveloso/soldim.git/internal/services/search"
	"github.com/claudineyveloso/soldim.git/internal/services/user"
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
	generatetoken.RegisterRoutes(r)
	product.RegisterRoutes(r)
	salechannel.RegisterRoutes(r)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)
	//
	searchStore := search.NewStore(s.db)
	searchHandler := search.NewHandler(searchStore)
	searchHandler.RegisterRoutes(r)
	//
	// ownerStore := owner.NewStore(s.db)
	// personStore := person.NewStore(s.db)
	// addressStore := address.NewStore(s.db)
	// ownerHandler := owner.NewHandler(ownerStore, personStore, addressStore)
	// ownerHandler.RegisterRoutes(r)
	//
	// customerStore := customer.NewStore(s.db)
	// personStore = person.NewStore(s.db)
	// addressStore = address.NewStore(s.db)
	// customerHandler := customer.NewHandler(customerStore, personStore, addressStore)
	// customerHandler.RegisterRoutes(r)
	//
	// typeServiceStore := typeservice.NewStore(s.db)
	// typeServiceHandler := typeservice.NewHandler(typeServiceStore)
	// typeServiceHandler.RegisterRoutes(r)
	//
	// intervalStore := interval.NewStore(s.db)
	// intervalHandler := interval.NewHandler(intervalStore)
	// intervalHandler.RegisterRoutes(r)
	//
	// attendanceStore := attendance.NewStore(s.db)
	// attendanceHandler := attendance.NewHandler(attendanceStore)
	// attendanceHandler.RegisterRoutes(r)
	//
	// insuranceStore := insurance.NewStore(s.db)
	// insuranceHandler := insurance.NewHandler(insuranceStore)
	// insuranceHandler.RegisterRoutes(r)
	//
	fmt.Println("Server started on http://localhost:8080")
	// return http.ListenAndServe("localhost:8080", r)
	return http.ListenAndServe("localhost:8080",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(r))
}
