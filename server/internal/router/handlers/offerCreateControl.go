package handler

import (
	"fmt"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func (h *Handlers) CreateOffer(IdStatus int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}
			requiredFields := []string{
				"price", //description is not requid
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}
			var offer models.Offer
			var err error
			offer.Price, err = strconv.ParseFloat(r.PostForm.Get("price"), 64)
			if err != nil {
				h.handleError(w, "Bad price", fmt.Errorf("failed price parsing"), 400)
				return
			}
			offer.Description = r.PostForm.Get("description")
			offer.DateTimeCreate = time.Now()

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

			idtender, err := h.getIDFromRequest(r)
			if err != nil {
				h.handleError(w, "Failed parse tender id", err, 500)
				return
			}
			tender, err := h.Repo.Tenders().GetByID(r.Context(), idtender)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
				return
			}
			if tender.IdCompany == user.IdCompany {
				h.handleError(w, "Нельзя создать коммерческое предложение на свой же тендер", fmt.Errorf("You cannot create a commercial offer for your own tender"), 409)
				return
			}
			offer.IdTender = tender.ID
			offer.IdCompany = user.IdCompany
			offer.IdStatus = IdStatus //Активный, 1 - Черновик todo в случае с кнопкой сохранить как черновик статус должен быть 1
			//validation
			err = offer.Validate()
			if err != nil {
				h.handleError(w, "Failed validate offer", err, 500)
				return
			}
			//parsing files
			files := r.MultipartForm.File["files"]

			if len(files) == 0 {
				h.handleError(w, "No files uploaded", fmt.Errorf("at least one file is required"), 400)
				return
			}
			for _, file := range files {
				err = validateTenderFile(file)
				if err != nil {
					h.handleError(w, err.Error(), err, 400)
				}
			}

			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Bad transaction", fmt.Errorf("Failed begin transaction: %v", err), 500)
				return
			}
			defer tx.Rollback()
			var success bool

			err = tx.Offer().Create(r.Context(), &offer)
			if err != nil {
				h.handleError(w, "Failed offer creation", fmt.Errorf("Failed offer creation in transaction: %v", err), 500)
				return
			}
			var log models.Log
			log.IdUser = user.ID
			log.IdEntity = 3
			log.IdType = 1
			log.DateTime = time.Now()

			err = tx.Log().Create(r.Context(), &log).Tender().Create(r.Context(), offer.ID)
			if err != nil {
				h.handleError(w, "Failed offer log creation", fmt.Errorf("Failed offer log creation in transaction: %v", err), 500)
				return
			}

			var path string
			path = fmt.Sprintf("./documents/offer%d", offer.ID)

			err = os.Mkdir(path, 0755)
			if err != nil {
				h.handleError(w, "Failed safed files", fmt.Errorf("Failed create mkdir: %v", err), 500)
				return
			}
			defer func() {
				if !success {
					tx.Rollback()
				} else {
					if err := tx.Commit(); err != nil {
						success = false
						errF := os.RemoveAll(path)
						if errF != nil {
							h.Logger.Errorf("Error of delete offers mkdir:%s;Error:%v", path, errF)
						}

						h.handleError(w, "Bad transaction", fmt.Errorf("Failed commit transaction: %v", err), 500)
						return
					}
				}
			}()
			for _, file := range files {
				var doc models.Doc
				doc.Name = file.Filename
				doc.FileName = path + "/" + file.Filename

				src, err := file.Open()
				if err != nil {
					h.handleError(w, "Failed to open uploaded file", fmt.Errorf("Failed open file: %v", err), 500)
					return
				}
				defer src.Close()

				f, err := os.Create(doc.FileName)
				if err != nil {
					h.handleError(w, "Failed to create file", fmt.Errorf("Failed create file: %v", err), 500)
					return
				}
				defer f.Close()

				_, err = io.Copy(f, src)
				if err != nil {
					h.handleError(w, "Failed to save file", fmt.Errorf("Failed save file: %v", err), 500)
					return
				}

				err = tx.Doc().Create(r.Context(), &doc).Offer().Create(r.Context(), offer.ID)
				if err != nil {
					os.Remove(doc.FileName)
					h.handleError(w, "Failed documents safed", fmt.Errorf("Failed doc creation in transaction: %v", err), 500)
					return
				}
				var log models.Log
				log.IdUser = user.ID
				log.IdEntity = 4
				log.IdType = 1
				log.DateTime = time.Now()
				err = tx.Log().Create(r.Context(), &log).Doc().Create(r.Context(), doc.ID)
				if err != nil {
					h.handleError(w, "Failed doc log creation", fmt.Errorf("Failed doc log creation in transaction: %v", err), 500)
					return
				}
			}

			success = true
			w.WriteHeader(201)
		}
	}
}
func (h *Handlers) GetCreateOfferWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/offer_create.html")
		if err != nil {
			h.handleError(w, "Failed tendercreate load:", err, 500)
			return
		}
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

		idtender, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Failed parse tender id", err, 500)
			return
		}
		tender, err := h.Repo.Tenders().GetByID(r.Context(), idtender)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		if tender.IdCompany == user.IdCompany {
			h.handleError(w, "Нельзя создать коммерческое предложение на свой же тендер", fmt.Errorf("You cannot create a commercial offer for your own tender"), 409)
			return
		}

		type Data struct {
			TenderID int
			Link     string
		}

		err = tmpl.Execute(w, &Data{Link: fmt.Sprintf("/main/tenders/%d", id), TenderID: tender.ID})
		if err != nil {
			h.handleError(w, "Failed tendercreate render:", err, 500)
			return
		}
	}
}
