package handler

import (
	"archive/zip"
	"fmt"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetTenderDocuments() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}
		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get docs:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Tender().GetAll(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get doc-links:", err, 500)
			return
		}

		var tenderdocs []models.Doc
		for _, link := range doclinks {
			for _, doc := range docs {
				if link.DocID == doc.ID {
					tenderdocs = append(tenderdocs, doc)
					break
				}
			}
		}
		if len(tenderdocs) > 0 {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename=\"documents.zip\"")
			zipWriter := zip.NewWriter(w)
			defer zipWriter.Close()

			for _, doc := range tenderdocs {
				file, err := os.Open(doc.FileName)
				if err != nil {
					continue
				}

				zipFile, err := zipWriter.Create(doc.Name)
				if err != nil {
					file.Close()
					continue
				}

				_, err = io.Copy(zipFile, file)
				file.Close()
				if err != nil {
					continue
				}
			}

		} else {
			h.handleError(w, "Tender not had documents", nil, 500)
		}
	}
}
func (h *Handlers) GetTenderDocumentById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}
		doc, err := h.Repo.Doc().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get doc:", err, 500)
			return
		}
		file, err := os.Open(doc.FileName)
		if err != nil {
			h.handleError(w, "Failed open file:", err, 500)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			h.handleError(w, "Failed get file info:", err, 500)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.Name))

		http.ServeContent(w, r, doc.Name, stat.ModTime(), file)
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
		if user.IdCompany != tender.IdCompany && tender.IdStatus == 1 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by tender"), 409)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_view.html")
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
				if offer.IdStatus == 1 || offer.IdStatus == 3{
					continue
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
				offersview = append(offersview, offerview)
			}
		}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tender           *views.TenderView
			Offers           []views.OfferView
		}
		err = tmpl.Execute(w, &data{Tender: tenview, LoginForm: LoginForm, RegistrationForm: RegistrationForm,Offers: offersview})
		if err != nil {
			h.handleError(w, "Failed tender render:", err, 500)
			return
		}
	}
}
