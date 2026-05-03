package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetMyOfferWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}

		offer, err := h.Repo.Offer().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get offer:", err, 500)
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
		if user.IdCompany != offer.IdCompany && offer.IdStatus == 1 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by offer"), 409)
			return
		}
		
		tmpl, err := template.ParseFiles("./client/pages/my_offer_view.html")
		if err != nil {
			h.handleError(w, "Failed offer load:", err, 500)
			return
		}
		company, err := h.Repo.Company().GetByID(r.Context(), offer.IdCompany)
		if err != nil {
			h.handleError(w, "Failed get company:", err, 500)
			return
		}
		status, err := h.Repo.Status().GetByID(r.Context(), offer.IdStatus)
		if err != nil {
			h.handleError(w, "Failed get status:", err, 500)
			return
		}

		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get files:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Offer().GetAll(r.Context(), offer.ID)
		if err != nil {
			h.handleError(w, "Failed get files-links:", err, 500)
			return
		}
		var files []models.Doc
		for _, link := range doclinks {
			for _, d := range docs {
				if d.ID == link.DocID {
					files = append(files, d)
					break
				} else {
					continue
				}
			}
		}
		var offersview *views.OfferView
		offersview = &views.OfferView{ID: offer.ID, Description: offer.Description, DateTimeCreate: GetDateString(offer.DateTimeCreate),
			Price: offer.Price, Company: company.Name, Status: status.Name, Files: files, IdTender: offer.IdTender,IdCompany: offer.IdCompany,}
		type data struct {
			Offer *views.OfferView
		}
		err = tmpl.Execute(w, &data{Offer: offersview})
		if err != nil {
			h.handleError(w, "Failed offer render:", err, 500)
			return
		}
	}
}
