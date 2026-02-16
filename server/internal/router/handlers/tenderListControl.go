package handler

import (
	"html/template"
	"main/internal/repositories/models"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetTendersListwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders,err:=h.Repo.Tenders().GetAll(r.Context())
		if err!=nil{
			h.handleError(w, "Failed get tenders:", err, 500)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_list.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		type data struct{
			Tenders []models.Tender
		}
		err = tmpl.Execute(w, &data{Tenders: tenders})
		if err != nil {
			h.handleError(w, "Failed profil render:", err, 500)
			return
		}
	}
}
