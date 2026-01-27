package router

import (
	handler "main/internal/router/handlers"
	"main/internal/router/handlers/jwt"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewRouter(h *handler.Handlers, db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()
	// Инициализация
	jwtService := jwt.NewService(db)

	r.HandleFunc("/auth/login", jwtService.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", jwtService.RefreshHandler).Methods("POST")
	r.HandleFunc("/auth/revoke", jwtService.RevokeHandler).Methods("POST")
	r.HandleFunc("/users", h.CreateUser()).Methods("POST")

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(jwtService.Middleware)
	// Users routes
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", h.GetUsers()).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", h.GetUserByID()).Methods("GET")
	userRouter.HandleFunc("/{login}", h.GetUserByLogin()).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", h.UpdateUser()).Methods("PUT")
	userRouter.HandleFunc("/{id:[0-9]+}", h.DeleteUser()).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", h.HealthCheck()).Methods("GET")

	return r
}
