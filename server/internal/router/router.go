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
	mainRouter.HandleFunc("/tenders", h.GetTendersListwindow()).Methods("GET")
	mainRouter.HandleFunc("/tenders/{id}", h.GetTenderwindow()).Methods("GET")
	mainRouter.HandleFunc("/registration", h.Registration()).Methods("POST")
	//authentification
	r.HandleFunc("/auth/login", jwtService.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", jwtService.RefreshHandler).Methods("GET")
	r.HandleFunc("/auth/revoke", jwtService.RevokeHandler).Methods("GET")
	//protected
	prRouter := r.PathPrefix("/protected").Subrouter()
	prRouter.Use(jwtService.Middleware)

	prRouter.HandleFunc("/lk", h.GetProfilwindow()).Methods("GET")
	prRouter.HandleFunc("/lk/edit", h.EditingLK()).Methods("POST")

	//---------------------------------------------------------------------------------------------\\
	// Health check
	r.HandleFunc("/health", h.HealthCheck()).Methods("GET")
	//testing
	r.HandleFunc("/users", h.GetUsers()).Methods("GET")
	r.HandleFunc("/logs", h.GetLogs()).Methods("GET")
	r.HandleFunc("/tenders", h.GetTenders()).Methods("GET")
	r.HandleFunc("/companies", h.GetCompanies()).Methods("GET")
	return r
}
