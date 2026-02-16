package handler

import (
	"encoding/json"
	"main/internal/repositories"
	"main/internal/repositories/models"
	"net/http"

	_ "github.com/lib/pq"
)

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
func (h *Handlers) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err,500)
			return
		}

		user, err := h.Repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed to get user", err,500)
			return
		}
		jsonResponse(w, user)
	}
}
func (h *Handlers) GetUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := h.getLoginFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid Login", err,500)
			return
		}

		user, err := h.Repo.Users().GetByLogin(r.Context(), login)
		if err != nil {
			h.handleError(w, "Failed to get user", err,500)
			return
		}
		jsonResponse(w, user)
	}
}
func (h *Handlers) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			h.handleError(w, "Invalid request body", err,500)
			return
		}

		if err := h.executeInTransaction(r, func(tx repositories.Transaction) error {
			return tx.Users().Create(r.Context(), &user)
		}); err != nil {
			return
		}

		jsonResponse(w, user, http.StatusCreated)
	}
}

func (h *Handlers) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			h.handleError(w, "Invalid request body", err,500)
			return
		}

		if err := h.executeInTransaction(r, func(tx repositories.Transaction) error {
			return tx.Users().Update(r.Context(), &user)
		}); err != nil {
			return
		}

		jsonResponse(w, user)
	}
}

func (h *Handlers) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err,500)
			return
		}

		if err := h.executeInTransaction( r, func(tx repositories.Transaction) error {
			return tx.Users().Delete(r.Context(), id)
		}); err != nil {
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
