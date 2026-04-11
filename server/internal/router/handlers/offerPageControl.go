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

func (h *Handlers) GetOfferDocuments() func(w http.ResponseWriter, r *http.Request) {
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
		doclinks, err := h.Repo.LinkerDoc(0).Offer().GetAll(r.Context(), id)
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
func (h *Handlers) GetOfferDocumentById() func(w http.ResponseWriter, r *http.Request) {
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
func (h *Handlers) GetOfferWindow() func(w http.ResponseWriter, r *http.Request) {
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

		if offer.IdStatus != 2 && offer.IdStatus != 6 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, bad status"), 409)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/offer_view.html")
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
			Price: offer.Price, Company: company.Name, Status: status.Name, Files: files, IdTender: offer.IdTender}
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
