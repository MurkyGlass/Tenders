package handler

import (
	"fmt"
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetMyAwaitTendersListwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders, err := h.Repo.Tenders().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get tenders:", err, 500)
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
		var DistMap map[int]string = make(map[int]string)
		for _, dist := range districts {
			DistMap[dist.ID] = dist.Name
		}
		var StMap map[int]string = make(map[int]string)
		for _, st := range statuses {
			StMap[st.ID] = st.Name
		}

		// сортировка фильтрации ограничения итд
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
		company, err := h.Repo.Company().GetByID(r.Context(), user.IdCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		//
		var TenderViews []views.TenderView
		for _, tender := range tenders {
			//  не черновик и не активный
			if tender.IdStatus != 6 {
				continue
			}
			if tender.IdCompany != company.ID {
				continue
			}

			TenderViews = append(TenderViews, views.TenderView{ID: tender.ID, Name: tender.Name,
				Description: tender.Description, DateTimeStart: GetDateString(tender.DateTimeStart),
				DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: company.Name, Status: StMap[tender.IdStatus],
				District: DistMap[tender.IdDistrict],IdCompany: tender.IdCompany,})
		}

		tmpl, err := template.ParseFiles("./client/pages/awaiting_tender_list.html")
		if err != nil {
			h.handleError(w, "Failed my await tenders load:", err, 500)
			return
		}

		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tenders          []views.TenderView
		}
		err = tmpl.Execute(w, &data{Tenders: TenderViews, LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed await tenderlist render:", err, 500)
			return
		}
	}
}
