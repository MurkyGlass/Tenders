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
	jwtService := jwt.NewService(db,h.Logger)
	//client static
	fs := http.FileServer(http.Dir("./client/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//main page
	mainRouter:=r.PathPrefix("/main").Subrouter()
	mainRouter.HandleFunc("",h.GetMainwindow()).Methods("GET")
	mainRouter.HandleFunc("/registration", h.Registration()).Methods("POST")
	//authentification
	r.HandleFunc("/auth/login", jwtService.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", jwtService.RefreshHandler).Methods("GET")
	r.HandleFunc("/auth/revoke", jwtService.RevokeHandler).Methods("GET")
	//protected
	apiRouter := r.PathPrefix("/protected").Subrouter()
	apiRouter.Use(jwtService.Middleware)
	
	apiRouter.HandleFunc("/lk",h.GetProfilwindow()).Methods("GET")

	// Health check
	r.HandleFunc("/health", h.HealthCheck()).Methods("GET")
	//testing
	r.HandleFunc("/users",h.GetUsers()).Methods("GET")

	return r
}
