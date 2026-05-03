package handler

import (
	"archive/zip"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetCompanyDocuments() func(w http.ResponseWriter, r *http.Request) {
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
		doclinks, err := h.Repo.LinkerDoc(0).Company().GetAll(r.Context(), id)
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

func (h *Handlers) GetCompanyWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}

		company, err := h.Repo.Company().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get tender:", err, 500)
			return
		}

		tmpl, err := template.ParseFiles("./client/pages/company_view.html")
		if err != nil {
			h.handleError(w, "Failed tender load:", err, 500)
			return
		}

		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get files:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Company().GetAll(r.Context(), company.ID)
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

		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Company          models.Company
			Files            []models.Doc
		}
		err = tmpl.Execute(w, &data{Company: *company, Files: files, LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed tender render:", err, 500)
			return
		}
	}
}
