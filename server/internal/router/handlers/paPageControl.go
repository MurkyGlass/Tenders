package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetProfilwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
		}
		user, err := h.repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
		}
		company, err := h.repo.Company().GetByID(r.Context(), user.IdCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
		}

		tmpl, err := template.ParseFiles("./client/pages/personal_account.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		type data struct {
			User    models.User
			Company models.Company
		}
		d := &data{User: *user, Company: *company}
		err = tmpl.Execute(w, d)
		if err != nil {
			h.handleError(w, "Failed profil render:", err, 500)
			return
		}
	}
}
