package handler

import (
	"html/template"
	"main/internal/repositories/models"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetCreateTenderWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/tender_create.html")
		if err != nil {
			h.handleError(w, "Failed tendercreate load:", err, 500)
			return
		}
		dist, err := h.Repo.District().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get districts:", err, 500)
			return
		}
		type Data struct {
			District []models.District
		}

		err = tmpl.Execute(w, &Data{District: dist})
		if err != nil {
			h.handleError(w, "Failed tendercreate render:", err, 500)
			return
		}
	}
}
