package handler

import (
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) PStext() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/PS.txt")
		if err != nil {
			h.handleError(w, "Failed help load:", err, 500)
			return
		}
		type Data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
		}
		err = tmpl.Execute(w, &Data{LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed help render:", err, 500)
			return
		}
	}
}
