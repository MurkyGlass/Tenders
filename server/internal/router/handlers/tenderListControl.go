package handler

import (
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetTendersListwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders, err := h.Repo.Tenders().GetAll(r.Context())

		if err != nil {
			h.handleError(w, "Failed get tenders:", err, 500)
			return
		}
		companies, err := h.Repo.Company().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get companies:", err, 500)
			return
		}
		statuses, err := h.Repo.Status().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get statuses:", err, 500)
			return
		}
		districts, err := h.Repo.District().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get districts:", err, 500)
			return
		}

		var CompMap map[int]string = make(map[int]string)
		for _, comp := range companies {
			CompMap[comp.ID] = comp.Name
		}
		var DistMap map[int]string = make(map[int]string)
		for _, dist := range districts {
			DistMap[dist.ID] = dist.Name
		}
		var StMap map[int]string = make(map[int]string)
		for _, st := range statuses {
			StMap[st.ID] = st.Name
		}
		// сортировка фильтрации ограничения итд
		var TenderViews []views.TenderView
		for _, tender := range tenders {
			// черновик, завершен(отобр. в отдельном списке)
			if tender.IdStatus == 1 || tender.IdStatus == 3 {
				continue
			}
			TenderViews = append(TenderViews, views.TenderView{ID: tender.ID, Name: tender.Name,
				Description: tender.Description, DateTimeStart: tender.DateTimeStart,
				DateTimeEnd: tender.DateTimeEnd, Company: CompMap[tender.IdCompany], Status: StMap[tender.IdStatus],
				District: DistMap[tender.IdDistrict]})
		}

		tmpl, err := template.ParseFiles("./client/pages/tender_list.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tenders          []views.TenderView
		}
		err = tmpl.Execute(w, &data{Tenders: TenderViews, LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed tenderlist render:", err, 500)
			return
		}
	}
}
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
			DateTimeStart: tender.DateTimeStart, DateTimeEnd: tender.DateTimeEnd, Company: company.Name,
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
