package handler

import (
	"fmt"
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetOfferListWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/my_offer_list.html")
		if err != nil {
			h.handleError(w, "Failed offerlist load:", err, 500)
			return
		}
		//pesronal data
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			return
		}
		user, err := h.Repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}

		var offersview []views.OfferView
		offers, err := h.Repo.Offer().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get offers", err, 500)
			return
		}
		for _, offer := range offers {
			if offer.IdCompany == user.IdCompany {
				var offerview views.OfferView
				offerview.ID = offer.ID
				offerview.Description = offer.Description
				offerview.Price = offer.Price
				offerview.DateTimeCreate = GetDateString(offer.DateTimeCreate)
				comp, err := h.Repo.Company().GetByID(r.Context(), offer.IdCompany)
				if err != nil {
					h.handleError(w, "Failed get company by offer", err, 500)
					return
				}
				offerview.Company = comp.Name
				st, err := h.Repo.Status().GetByID(r.Context(), offer.IdStatus)
				offerview.Status = st.Name
				offerview.IdCompany = offer.IdCompany
				offersview = append(offersview, offerview)
			}
		}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Offers           []views.OfferView
		}
		err = tmpl.Execute(w, &data{LoginForm: LoginForm, RegistrationForm: RegistrationForm, Offers: offersview})
		if err != nil {
			h.handleError(w, "Failed offerList render:", err, 500)
			return
		}
	}
}
