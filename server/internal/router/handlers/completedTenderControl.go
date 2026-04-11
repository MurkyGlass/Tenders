package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func (h *Handlers) CompletTender() func(w http.ResponseWriter, r *http.Request) {
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
		if user.IdCompany != tender.IdCompany {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by tender"), 409)
			return
		}
		idoffer, err := h.getIDOfferFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}
		offer, err := h.Repo.Offer().GetByID(r.Context(), idoffer)
		if err != nil {
			h.handleError(w, "Failed get tender:", err, 500)
			return
		}
		if tender.IdStatus != 6 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, status"), 409)
			return
		}
		if offer.IdStatus != 2 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, status"), 409)
			return
		}
		offer.IdStatus = 5
		tender.IdStatus = 3
		tx, err := h.Repo.BeginTx(r.Context())
		if err != nil {
			h.handleError(w, "Bad transaction", fmt.Errorf("Failed begin transaction: %v", err), 500)
			return
		}
		defer tx.Rollback()
		err = tx.Tenders().Update(r.Context(),tender)
		if err != nil {
			h.handleError(w, "Bad update", fmt.Errorf("Failed update tender in transaction: %v", err), 500)
			return
		}
		var log models.Log
		log.DateTime = time.Now()
		log.IdType = 2
		log.IdEntity = 2
		log.IdUser = user.ID
		err = tx.Log().Create(r.Context(),&log).Tender().Create(r.Context(),tender.ID)
		if err != nil {
			h.handleError(w, "Bad loging", fmt.Errorf("Failed update tender logs in transaction: %v", err), 500)
			return
		}

		err = tx.Offer().Update(r.Context(),offer)
		if err != nil {
			h.handleError(w, "Bad update", fmt.Errorf("Failed update offer in transaction: %v", err), 500)
			return
		}
		log.ID = 0
		log.DateTime = time.Now()
		log.IdType = 2
		log.IdEntity = 3
		log.IdUser = user.ID
		err = tx.Log().Create(r.Context(),&log).Offer().Create(r.Context(),offer.ID)
		if err != nil {
			h.handleError(w, "Bad loging", fmt.Errorf("Failed update offer logs in transaction: %v", err), 500)
			return
		}
		//protocol create
		
		//
		err = tx.Commit()
		if err != nil {
			h.handleError(w, "Bad commit", fmt.Errorf("Failed commit in transaction: %v", err), 500)
			return
		}
		w.WriteHeader(201)
	}
}
func (h *Handlers) GetCompletedTenderwindow() func(w http.ResponseWriter, r *http.Request) {
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
		if user.IdCompany != tender.IdCompany {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by tender"), 409)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/completed_tender.html")
		if err != nil {
			h.handleError(w, "Failed tender load:", err, 500)
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
		categories, err := h.Repo.Category().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get categories:", err, 500)
			return
		}
		categoriesLinks, err := h.Repo.LinkerTCategory(tender.ID).GetAllByTender(r.Context())
		if err != nil {
			h.handleError(w, "Failed get c-links:", err, 500)
			return
		}
		var catlist []string
		for _, link := range categoriesLinks {
			for _, cat := range categories {
				if link.IdCategory == cat.ID {
					catlist = append(catlist, cat.Name)
					break
				} else {
					continue
				}
			}
		}
		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get files:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Tender().GetAll(r.Context(), tender.ID)
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
		var tenview *views.TenderView
		tenview = &views.TenderView{ID: tender.ID, Name: tender.Name, Description: tender.Description,
			DateTimeStart: GetDateString(tender.DateTimeStart), DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: company.Name,
			Status: status.Name, District: district.Name, Categories: catlist, Files: files}
		var offersview []views.OfferView
		offers, err := h.Repo.Offer().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get offers", err, 500)
			return
		}
		for _, offer := range offers {
			if offer.IdTender == tender.ID {
				if offer.IdStatus != 2 {
					continue
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
				offerview.Files = files
				offersview = append(offersview, offerview)
			}
		}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tender           *views.TenderView
			Offers           []views.OfferView
		}
		err = tmpl.Execute(w, &data{Tender: tenview, LoginForm: LoginForm, RegistrationForm: RegistrationForm, Offers: offersview})
		if err != nil {
			h.handleError(w, "Failed tender-completed render:", err, 500)
			return
		}
	}
}
