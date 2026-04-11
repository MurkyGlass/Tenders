package handler

import (
	"fmt"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func (h *Handlers) EditOffer(IdStatus int) func(w http.ResponseWriter, r *http.Request) {
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
			id, err := h.getIDFromRequest(r)
			if err != nil {
				h.handleError(w, "Failed id parsing", err, 500)
				return
			}
			offer, err := h.Repo.Offer().GetByID(r.Context(), id)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
				return
			}
			if offer.IdStatus != 1 && offer.IdStatus != 2 {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, status"), 409)
				return
			}
			offer.Price, err = strconv.ParseFloat(r.PostForm.Get("price"), 64)
			if err != nil {
				h.handleError(w, "Bad price", fmt.Errorf("failed price parsing"), 400)
				return
			}
			offer.Description = r.PostForm.Get("description")

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
			if user.IdCompany != offer.IdCompany {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by offer"), 409)
				return
			}

			offer.IdStatus = IdStatus //Активный, 1 - Черновик todo в случае с кнопкой сохранить как черновик статус должен быть 1
			//validation
			err = offer.Validate()
			if err != nil {
				h.handleError(w, "Failed validate offer", err, 500)
				return
			}
			//parsing files
			files := r.MultipartForm.File["files"]

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

			err = tx.Offer().Update(r.Context(), offer)
			if err != nil {
				h.handleError(w, "Failed offer creation", fmt.Errorf("Failed offer creation in transaction: %v", err), 500)
				return
			}
			var log models.Log
			log.IdUser = user.ID
			log.IdEntity = 3
			log.IdType = 2
			log.DateTime = time.Now()

			err = tx.Log().Create(r.Context(), &log).Offer().Create(r.Context(), offer.ID)
			if err != nil {
				h.handleError(w, "Failed offer log creation", fmt.Errorf("Failed offer log creation in transaction: %v", err), 500)
				return
			}

			var path string
			path = fmt.Sprintf("./documents/offer%d", offer.ID)

			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.Mkdir(path, 0755)
				if err != nil {
					h.handleError(w, "Failed create directory", fmt.Errorf("mkdir failed: %v", err), 500)
					return
				}
			}

			defer func() {
				if !success {
					tx.Rollback()
				} else {
					if err := tx.Commit(); err != nil {
						success = false
						errF := os.RemoveAll(path)
						if errF != nil {
							h.Logger.Errorf("Error of delete tenders mkdir:%s;Error:%v", path, errF)
						}

						h.handleError(w, "Bad transaction", fmt.Errorf("Failed commit transaction: %v", err), 500)
						return
					}
				}
			}()
			//files aus deleting
			IdsFileStr := r.PostForm.Get("oldfiles")
			var IdsFiles []int
			if IdsFileStr != "" {
				IDsStrArr := r.PostForm["oldfiles"]

				for _, str := range IDsStrArr {
					id, err := strconv.Atoi(str)
					if err != nil {
						h.handleError(w, "Invalid Parsing category id by Filter", err, 400)
						return
					}
					IdsFiles = append(IdsFiles, id)
				}
			}
			Docs, err := h.Repo.Doc().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed get docs", fmt.Errorf("Failed get docs: %v", err), 500)
				return
			}
			doclinks, err := h.Repo.LinkerDoc(0).Offer().GetAll(r.Context(), offer.ID)
			var OldDocs []models.Doc
			for _, l := range doclinks {
				for _, f := range Docs {
					if f.ID == l.DocID {
						OldDocs = append(OldDocs, f)
						break
					}
				}
			}
			var deliting []models.Doc
			for _, l := range IdsFiles {
				for _, f := range Docs {
					if f.ID == l {
						deliting = append(deliting, f)
						break
					}
				}
			}
			for _, d := range deliting {
				err = os.Remove(d.FileName)
				if err != nil {
					h.Logger.Errorf("Failed delete file:%v", err)
				}
				err = tx.LinkerLog(0).Doc().DeleteAll(r.Context(), d.ID)
				if err != nil {
					h.handleError(w, "Failed delete doc-log", fmt.Errorf("Failed delete doc-log: %v", err), 500)
					return
				}
				err = tx.LinkerDoc(d.ID).Offer().Delete(r.Context(), offer.ID)
				if err != nil {
					h.handleError(w, "Failed delete doc-tenders", fmt.Errorf("Failed delete doc-tenders: %v", err), 500)
					return
				}
				err = tx.Doc().Delete(r.Context(), d.ID)
				if err != nil {
					h.handleError(w, "Failed delete doc", fmt.Errorf("Failed delete doc: %v", err), 500)
					return
				}
			}

			for _, file := range files {
				flag := false
				for _, d := range OldDocs {
					if d.Name == file.Filename {
						flag = true
						for _, d := range deliting {
							if d.Name == file.Filename {
								flag = false
							}
						}
					}
				}
				if flag {
					continue
				}
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
func (h *Handlers) GetEditOfferWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/offer_edit.html")
		if err != nil {
			h.handleError(w, "Failed offer edit load:", err, 500)
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
		idoffer, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Failed parse tender id", err, 500)
			return
		}
		offer, err := h.Repo.Offer().GetByID(r.Context(), idoffer)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		if user.IdCompany != offer.IdCompany {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by offer"), 409)
			return
		}
		if offer.IdStatus != 1 && offer.IdStatus != 2 {
			h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, status"), 409)
			return
		}
		status, err := h.Repo.Status().GetByID(r.Context(), offer.IdStatus)
		if err != nil {
			h.handleError(w, "Failed get status", err, 500)
			return
		}
		company, err := h.Repo.Company().GetByID(r.Context(), offer.IdCompany)
		if err != nil {
			h.handleError(w, "Failed get company", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Offer().GetAll(r.Context(), offer.ID)
		if err != nil {
			h.handleError(w, "Failed get doclink", err, 500)
			return
		}
		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get docs", err, 500)
			return
		}
		var files []models.Doc
		for _, link := range doclinks {
			for _, d := range docs {
				if d.ID == link.DocID && link.OfferID == offer.ID {
					files = append(files, d)
					break
				}
			}
		}
		var view *views.OfferView
		view = &views.OfferView{
			Price: offer.Price, ID: offer.ID, IdTender: offer.IdTender, Description: offer.Description,
			DateTimeCreate: GetDateInputFormat(offer.DateTimeCreate), Status: status.Name, Company: company.Name, Files: files,
		}

		type Data struct {
			Offer *views.OfferView
		}

		err = tmpl.Execute(w, &Data{Offer: view})
		if err != nil {
			h.handleError(w, "Failed offeredit render:", err, 500)
			return
		}
	}
}
