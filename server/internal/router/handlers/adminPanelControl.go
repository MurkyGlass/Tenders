package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetAdminPanelWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/adminpanel.html")
		if err != nil {
			h.handleError(w, "Failed admin load:", err, 500)
			return
		}
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
		}
		user, err := h.Repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
		}
		type Data struct {
			User models.User
		}
		err = tmpl.Execute(w, &Data{User: *user})
		if err != nil {
			h.handleError(w, "Failed admin render:", err, 500)
			return
		}
	}
}
