package handler

import (
	"fmt"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func (h *Handlers) EditDraftTender(IdStatus int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), EditTenderPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}
			requiredFields := []string{
				"name", "id_district", "datetime_start", "datetime_end", //description is not requid
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
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
			if !tender.DateTimeEnd.After(time.Now()) {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict,by time end and Now"), 409)
				return
			}
			if tender.IdStatus != 1 {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict,status"), 409)
				return
			}
			//edit
			tender.Name = r.PostForm.Get("name")
			tender.Description = r.PostForm.Get("description")
			tender.IdDistrict, err = strconv.Atoi(r.PostForm.Get("id_district"))
			if err != nil {
				h.handleError(w, "Bad district", fmt.Errorf("failed district id parsing"), 400)
				return
			}
			const layout = "2006-01-02T15:04"
			tender.DateTimeStart, err = time.Parse(layout, r.PostForm.Get("datetime_start"))
			if err != nil {
				h.handleError(w, "Failed datetime parsing", fmt.Errorf("failed date parsing"), 400)
				return
			}
			tender.DateTimeEnd, err = time.Parse(layout, r.PostForm.Get("datetime_end"))
			if err != nil {
				h.handleError(w, "Failed datetime parsing", fmt.Errorf("failed date parsing"), 400)
				return
			}
			if !tender.DateTimeStart.Before(tender.DateTimeEnd) {
				h.handleError(w, "Date start must date end", fmt.Errorf("failed, date start must date end"), 400)
				return
			}
			if !tender.DateTimeStart.After(time.Now()) {
				h.handleError(w, "Date start must date now", fmt.Errorf("failed, date start must date now"), 400)
				return
			}

			id, ok := r.Context().Value("id_user").(int)
			if !ok {
				h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			}
			user, err := h.Repo.Users().GetByID(r.Context(), id)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
			}

			if user.IdCompany != tender.IdCompany {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by tender"), 409)
				return
			}
			tender.IdStatus = IdStatus //Черновик todo в случае с кнопкой сохранить как черновик статус должен быть 1

			//Validation
			err = tender.Validate()
			if err != nil {
				h.handleError(w, "tender validation failed", err, 500)
			}
			//parsing files//todo
			files := r.MultipartForm.File["files"]

			for _, file := range files {
				err = validateTenderFile(file)
				if err != nil {
					h.handleError(w, err.Error(), err, 400)
				}
			}

			IdsStr := r.PostForm.Get("category_ids")
			var IdsCateg []int
			if IdsStr != "" {
				IDsStrArr := r.PostForm["category_ids"]

				for _, str := range IDsStrArr {
					id, err := strconv.Atoi(str)
					if err != nil {
						h.handleError(w, "Invalid Parsing category id by Filter", err, 400)
						return
					}
					IdsCateg = append(IdsCateg, id)
				}
			}

			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Bad transaction", fmt.Errorf("Failed begin transaction: %v", err), 500)
				return
			}
			defer tx.Rollback()
			var success bool

			err = tx.Tenders().Update(r.Context(), tender)
			if err != nil {
				h.handleError(w, "Failed tender updating", fmt.Errorf("Failed tender updating in transaction: %v", err), 500)
				return
			}
			var log models.Log
			log.IdUser = user.ID
			log.IdEntity = 2
			log.IdType = 2
			log.DateTime = time.Now()

			err = tx.Log().Create(r.Context(), &log).Tender().Create(r.Context(), tender.ID)
			if err != nil {
				h.handleError(w, "Failed tender log creation", fmt.Errorf("Failed tender log creation in transaction: %v", err), 500)
				return
			}

			links, err := tx.CategoryLink().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed get categ links", fmt.Errorf("Failed get categ links in transaction: %v", err), 500)
				return
			}
			if len(IdsCateg) > 0 {
				var parents []int
				var f bool
				f = false
				for !f {
					f = true
					for _, i := range IdsCateg {
						for _, l := range links {
							if l.SecondID == i {
								if !slices.ContainsFunc(IdsCateg, func(e int) bool {
									return e == l.FirstID
								}) && !slices.ContainsFunc(parents, func(a int) bool {
									return a == l.FirstID
								}) {
									parents = append(parents, l.FirstID)
								}

								links = slices.DeleteFunc(links, func(link models.LinkView) bool {
									return link.FirstID == l.FirstID && link.SecondID == l.SecondID
								})
								break
							}
						}
					}
					if len(parents) > 0 {
						for _, p := range parents {
							IdsCateg = append(IdsCateg, p)
						}
						parents = nil
						parents = make([]int, 0)
						f = false
					}
				}
			}

			if len(IdsCateg) > 0 {
				for _, i := range IdsCateg {
					err = tx.LinkerTCategory(tender.ID).Delete(r.Context(), i)
					if err != nil {
						h.handleError(w, "Failed tender-category link deliting", fmt.Errorf("Failed tender-category link  deleting in transaction: %v", err), 500)
						return
					}
				}
				for _, i := range IdsCateg {
					err = tx.LinkerTCategory(tender.ID).Create(r.Context(), i)
					if err != nil {
						h.handleError(w, "Failed tender-category link creation", fmt.Errorf("Failed tender-category link  creation in transaction: %v", err), 500)
						return
					}
				}
			}
			//
			var path string
			path = fmt.Sprintf("./documents/tender%d", tender.ID)
			//добавить создания папки но с проверкой
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
			doclinks, err := h.Repo.LinkerDoc(0).Tender().GetAll(r.Context(), tender.ID)
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
				err = tx.LinkerDoc(d.ID).Tender().Delete(r.Context(), tender.ID)
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

				err = tx.Doc().Create(r.Context(), &doc).Tender().Create(r.Context(), tender.ID)
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
func (h *Handlers) EditTender() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), EditPublishTenderPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}
			requiredFields := []string{
				"id_district",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
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
			if !tender.DateTimeEnd.After(time.Now()) {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict,by time end and Now"), 409)
				return
			}
			if tender.IdStatus != 1 && tender.IdStatus != 2 {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict,status"), 409)
				return
			}
			//edit
			tender.Description = r.PostForm.Get("description")
			tender.IdDistrict, err = strconv.Atoi(r.PostForm.Get("id_district"))
			if err != nil {
				h.handleError(w, "Bad district", fmt.Errorf("failed district id parsing"), 400)
				return
			}

			id, ok := r.Context().Value("id_user").(int)
			if !ok {
				h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			}
			user, err := h.Repo.Users().GetByID(r.Context(), id)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
			}

			if user.IdCompany != tender.IdCompany {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by tender"), 409)
				return
			}
			tender.IdStatus = 2 //Активный, 1 - Черновик todo в случае с кнопкой сохранить как черновик статус должен быть 1

			//Validation
			err = tender.Validate()
			if err != nil {
				h.handleError(w, "tender validation failed", err, 500)
			}
			//parsing files//todo
			files := r.MultipartForm.File["files"]

			for _, file := range files {
				err = validateTenderFile(file)
				if err != nil {
					h.handleError(w, err.Error(), err, 400)
				}
			}

			IdsStr := r.PostForm.Get("category_ids")
			var IdsCateg []int
			if IdsStr != "" {
				IDsStrArr := r.PostForm["category_ids"]

				for _, str := range IDsStrArr {
					id, err := strconv.Atoi(str)
					if err != nil {
						h.handleError(w, "Invalid Parsing category id by Filter", err, 400)
						return
					}
					IdsCateg = append(IdsCateg, id)
				}
			}

			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Bad transaction", fmt.Errorf("Failed begin transaction: %v", err), 500)
				return
			}
			defer tx.Rollback()
			var success bool

			err = tx.Tenders().Update(r.Context(), tender)
			if err != nil {
				h.handleError(w, "Failed tender updating", fmt.Errorf("Failed tender updating in transaction: %v", err), 500)
				return
			}
			var log models.Log
			log.IdUser = user.ID
			log.IdEntity = 2
			log.IdType = 2
			log.DateTime = time.Now()

			err = tx.Log().Create(r.Context(), &log).Tender().Create(r.Context(), tender.ID)
			if err != nil {
				h.handleError(w, "Failed tender log creation", fmt.Errorf("Failed tender log creation in transaction: %v", err), 500)
				return
			}

			links, err := tx.CategoryLink().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed get categ links", fmt.Errorf("Failed get categ links in transaction: %v", err), 500)
				return
			}
			if len(IdsCateg) > 0 {
				var parents []int
				var f bool
				f = false
				for !f {
					f = true
					for _, i := range IdsCateg {
						for _, l := range links {
							if l.SecondID == i {
								if !slices.ContainsFunc(IdsCateg, func(e int) bool {
									return e == l.FirstID
								}) && !slices.ContainsFunc(parents, func(a int) bool {
									return a == l.FirstID
								}) {
									parents = append(parents, l.FirstID)
								}

								links = slices.DeleteFunc(links, func(link models.LinkView) bool {
									return link.FirstID == l.FirstID && link.SecondID == l.SecondID
								})
								break
							}
						}
					}
					if len(parents) > 0 {
						for _, p := range parents {
							IdsCateg = append(IdsCateg, p)
						}
						parents = nil
						parents = make([]int, 0)
						f = false
					}
				}
			}

			if len(IdsCateg) > 0 {
				for _, i := range IdsCateg {
					err = tx.LinkerTCategory(tender.ID).Delete(r.Context(), i)
					if err != nil {
						h.handleError(w, "Failed tender-category link deliting", fmt.Errorf("Failed tender-category link  deleting in transaction: %v", err), 500)
						return
					}
				}
				for _, i := range IdsCateg {
					err = tx.LinkerTCategory(tender.ID).Create(r.Context(), i)
					if err != nil {
						h.handleError(w, "Failed tender-category link creation", fmt.Errorf("Failed tender-category link  creation in transaction: %v", err), 500)
						return
					}
				}
			}
			//
			var path string
			path = fmt.Sprintf("./documents/tender%d", tender.ID)
			//добавить создания папки но с проверкой
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
			doclinks, err := h.Repo.LinkerDoc(0).Tender().GetAll(r.Context(), tender.ID)
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
				err = tx.LinkerDoc(d.ID).Tender().Delete(r.Context(), tender.ID)
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

				err = tx.Doc().Create(r.Context(), &doc).Tender().Create(r.Context(), tender.ID)
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
func BuildCategoryTreeByChecked(categories []models.Category, links []models.LinkView, selected map[int]bool) []views.CategoryView {
	var cm = make(map[int]models.Category)
	var lpk = make(map[int][]int) //parrent key
	var lck = make(map[int]int)   //child key

	for _, category := range categories {
		cm[category.ID] = category
	}
	for _, link := range links {
		lpk[link.FirstID] = append(lpk[link.FirstID], link.SecondID)
		lck[link.SecondID] = link.FirstID
	}

	var view []views.CategoryView
	var childs = make(map[int]models.Category)
	for _, c := range categories {
		_, f0 := lck[c.ID]
		if f0 { //если явл ребенком
			childs[c.ID] = c
			continue
		}
		_, f1 := lpk[c.ID]
		if f1 && !f0 { //если явл отцом и не ребенок
			view = append(view, views.CategoryView{Category: c, Childs: nil, Checked: selected[c.ID]})
			continue
		}
		view = append(view, views.CategoryView{Category: c, Childs: nil, Checked: selected[c.ID]})
	}

	var road []*views.CategoryView
	for i := range view {
		road = append(road, &view[i])
	}

	for len(childs) > 0 {

		for _, r := range road {
			for _, l := range links {
				if l.FirstID == r.Category.ID {
					r.Childs = append(r.Childs, views.CategoryView{Category: childs[l.SecondID], Childs: nil, Checked: selected[childs[l.SecondID].ID]})
					delete(childs, l.SecondID)
				}
			}
		}

		newRoad := []*views.CategoryView{}
		for _, r := range road {
			for i := range r.Childs {
				newRoad = append(newRoad, &r.Childs[i])
			}
		}
		road = newRoad
	}
	return view
}

// good
func (h *Handlers) GetEditTenderWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), EditTenderPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_edit.html")
		if err != nil {
			h.handleError(w, "Failed tenderedit load:", err, 500)
			return
		}
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
		if tender.IdStatus != 1 && tender.IdStatus != 2 {
			h.handleError(w, "Tender status not active or draft:", fmt.Errorf("Tender status not active or draft"), 409)
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

		categories, err := h.Repo.Category().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get categories:", err, 500)
			return
		}
		catlinks, err := h.Repo.CategoryLink().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get category links:", err, 500)
			return
		}
		tl, err := h.Repo.LinkerTCategory(tender.ID).GetAllByTender(r.Context())
		if err != nil {
			h.handleError(w, "Failed get category-tender links:", err, 500)
			return
		}
		var selected = make(map[int]bool)
		for _, l := range tl {
			selected[l.IdCategory] = true
		}
		cv := BuildCategoryTreeByChecked(categories, catlinks, selected) //edit by checked
		dist, err := h.Repo.District().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get districts:", err, 500)
			return
		}
		alldoc, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get docs:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Tender().GetAll(r.Context(), tender.ID)
		if err != nil {
			h.handleError(w, "Failed get doc-links:", err, 500)
			return
		}
		var docs []models.Doc
		for _, link := range doclinks {
			for _, doc := range alldoc {
				if doc.ID == link.DocID {
					docs = append(docs, doc)
				}
			}
		}
		type Data struct {
			District      []models.District
			CatView       []views.CategoryView
			Tender        *models.Tender
			OldFiles      []models.Doc
			DateTimeStart string
			DateTimeEnd   string
		}

		err = tmpl.Execute(w, &Data{District: dist, CatView: cv, Tender: tender, OldFiles: docs, DateTimeStart: GetDateInputFormat(tender.DateTimeStart),
			DateTimeEnd: GetDateInputFormat(tender.DateTimeEnd)})
		if err != nil {
			h.handleError(w, "Failed tenderedit render:", err, 500)
			return
		}
	}
}
func GetDateInputFormat(t time.Time) string {
	month := t.Month()
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	var monstr string
	if int(month) < 10 {
		monstr = "0" + strconv.Itoa(int(month))
	} else {
		monstr = strconv.Itoa(int(month))
	}
	var daystr string
	if day < 10 {
		daystr = "0" + strconv.Itoa(day)
	} else {
		daystr = strconv.Itoa(day)
	}
	var hstr string
	if hour < 10 {
		hstr = "0" + strconv.Itoa(hour)
	} else {
		hstr = strconv.Itoa(hour)
	}
	var mstr string
	if minute < 10 {
		mstr = "0" + strconv.Itoa(minute)
	} else {
		mstr = strconv.Itoa(minute)
	}
	return fmt.Sprintf("%d-%s-%sT%s:%s", t.Year(), monstr, daystr, hstr, mstr)

}
