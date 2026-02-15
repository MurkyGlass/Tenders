package handler

import (
	"encoding/json"
	"fmt"
	"main/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	Logger *logrus.Logger
	Repo   repositories.Repository
}

func NewHandlers(logger *logrus.Logger, rep repositories.Repository) *Handlers {
	return &Handlers{
		Logger: logger,
		Repo:   rep,
	}
}

func (h *Handlers) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// Логирование
func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// CORS
func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handlers) getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	return strconv.Atoi(idStr)
}

func (h *Handlers) getLoginFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	Str := vars["login"]
	l := len(Str)
	if l > 3 && l < 31 {
		return Str, nil
	}
	return Str, fmt.Errorf("Validation Error;%s; len in bait no is 3<l<31", Str)
}

// Runing in TX
func (h *Handlers) executeInTransaction( r *http.Request, fn func(tx repositories.Transaction) error) error {
	tx, err := h.Repo.BeginTx(r.Context())
	if err != nil {
		h.Logger.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		h.Logger.Errorf("Operation failed in TX: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		h.Logger.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (h *Handlers) handleError(w http.ResponseWriter, message string, err error, code int) {
	h.Logger.Errorf("%s: %v", message, err)
	http.Error(w, message, code)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	json.NewEncoder(w).Encode(data)
}
