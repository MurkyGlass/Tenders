package handler

import (
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetTenderwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}

		tender, err := h.Repo.Tenders().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get tender:", err, 500)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_view.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		company, err := h.Repo.Company().GetByID(r.Context(), tender.IdCompany)
		if err != nil {
			h.handleError(w, "Failed get company:", err, 500)
			return
		}
		status, err := h.Repo.Status().GetByID(r.Context(), tender.IdStatus)
		if err != nil {
			h.handleError(w, "Failed get status:", err, 500)
			return
		}
		district, err := h.Repo.District().GetByID(r.Context(), tender.IdDistrict)
		if err != nil {
			h.handleError(w, "Failed get district:", err, 500)
			return
		}
		var tenview *views.TenderView
		tenview = &views.TenderView{ID: tender.ID, Name: tender.Name, Description: tender.Description,
			DateTimeStart: GetDateString(tender.DateTimeStart), DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: company.Name,
			Status: status.Name, District: district.Name}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tender           *views.TenderView
		}
		err = tmpl.Execute(w, &data{Tender: tenview, LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed tender render:", err, 500)
			return
		}
	}
}
