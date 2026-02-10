package handler

import (
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetProfilwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/personal_account.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			h.handleError(w, "Failed profil render:", err, 500)
			return
		}
	}
}
