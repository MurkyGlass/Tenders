package handler

import (
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetMainwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/main.html")
		if err != nil {
			h.handleError(w, "Failed main load:", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			h.handleError(w, "Failed main render:", err)
			return
		}
	}
}
func (h *Handlers) GetProfilwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/main.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			h.handleError(w, "Failed profil render:", err)
			return
		}
	}
}