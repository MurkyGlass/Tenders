package handler

import (
	"net/http"

	_ "github.com/lib/pq"
)

/*
	func (h *Handlers) GetUsers() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			users, err := h.Repo.Users().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed to get users", err,500)
				return
			}
			jsonResponse(w, users)
		}
	}
*/
func (h *Handlers) GetLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := h.Repo.Log().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed to get logs", err, 500)
			return
		}
		jsonResponse(w, logs)
	}
}
func (h *Handlers) GetTenders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders, err := h.Repo.Tenders().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed to get tenders", err, 500)
			return
		}
		jsonResponse(w, tenders)
	}
}
func (h *Handlers) GetCompanies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := h.Repo.Company().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed to get companies", err, 500)
			return
		}
		jsonResponse(w, c)
	}
}
