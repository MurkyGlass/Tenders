package router

import (
	handler "main/internal/router/handlers"
	"main/internal/router/handlers/jwt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewRouter(h *handler.Handlers, db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()
	// Инициализация
	jwtService := jwt.NewService(h.Logger, h.Repo)
	//client static
	fs := http.FileServer(http.Dir("./client/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//main page
	r.HandleFunc("", h.GetMainwindow()).Methods("GET")
	mainRouter := r.PathPrefix("/main").Subrouter()
	mainRouter.HandleFunc("", h.GetMainwindow()).Methods("GET")
	mainRouter.HandleFunc("/tenders", h.GetTendersListwindow(nil)).Methods("GET")
	mainRouter.HandleFunc("/tenders", h.FilterParams()).Methods("POST")
	mainRouter.HandleFunc("/tenders/{id}", h.GetTenderwindow()).Methods("GET")
	mainRouter.HandleFunc("/registration", h.Registration()).Methods("POST")

	//authentification
	r.HandleFunc("/auth/login", jwtService.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", jwtService.RefreshHandler).Methods("GET")
	r.HandleFunc("/auth/revoke", jwtService.RevokeHandler).Methods("GET")
	//protected
	prRouter := r.PathPrefix("/protected").Subrouter()
	prRouter.Use(jwtService.Middleware)

	//tenders docs
	prRouter.HandleFunc("/tenders/documents/{id}", h.GetTenderDocumentById()).Methods("GET")
	prRouter.HandleFunc("/tenders/{id}/documents", h.GetTenderDocuments()).Methods("GET")
	//--------------------------------------------------------------------------------------
	//offer-create
	prRouter.HandleFunc("/tenders/{id}/offer/create",h.GetCreateOfferWindow()).Methods("GET")
	prRouter.HandleFunc("/tenders/{id}/offer/create",h.CreateOffer(2)).Methods("POST")
	prRouter.HandleFunc("/tenders/{id}/offer/create/draft",h.CreateOffer(1)).Methods("POST")
	//--------------------------------------------------------------------------------------
	//offers
	prRouter.HandleFunc("/offers/{id}",h.GetOfferWindow()).Methods("GET")
	//offers docs
	prRouter.HandleFunc("/offers/documents/{id}", h.GetOfferDocumentById()).Methods("GET")
	prRouter.HandleFunc("/offers/{id}/documents", h.GetOfferDocuments()).Methods("GET")
	//--------------------------------------------------------------------------------------
	lkRouter := prRouter.PathPrefix("/lk").Subrouter()
	lkRouter.HandleFunc("", h.GetProfilwindow()).Methods("GET")
	lkRouter.HandleFunc("/edit", h.EditingLK()).Methods("POST")

	tenderRouter := prRouter.PathPrefix("/tender").Subrouter()
	tenderRouter.HandleFunc("/create", h.GetCreateTenderWindow()).Methods("GET")
	tenderRouter.HandleFunc("/create", h.CreateTender(2)).Methods("POST")
	tenderRouter.HandleFunc("/create/draft", h.CreateTender(1)).Methods("POST")
	//---------------------------------------------------------------------------------------------\\
	// Health check
	r.HandleFunc("/health", h.HealthCheck()).Methods("GET")
	//testing
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/logs", h.GetLogs()).Methods("GET")
	r.HandleFunc("/tenders", h.GetTenders()).Methods("GET")
	r.HandleFunc("/tender/{id}/categories", h.GetTenderCategoryLinksbyId()).Methods("GET")
	r.HandleFunc("/tenders/categories", h.GetTenderCategoryLinks()).Methods("GET")
	r.HandleFunc("/companies", h.GetCompanies()).Methods("GET")
	return r
}
