package handler

import (
	"fmt"
	"html/template"
	"io"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	_ "github.com/lib/pq"
)

var allowedMimeTypes = map[string]string{
	// Документы
	"application/pdf":    "PDF документ",
	"application/msword": "Word документ (старый формат)",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "Word документ",
	"application/vnd.ms-excel": "Excel таблица (старый формат)",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "Excel таблица",

	// Архивы
	"application/zip":              "ZIP архив",
	"application/x-zip-compressed": "ZIP архив",

	// XML
	"application/xml":       "XML файл",
	"text/xml":              "XML файл",
	"application/rss+xml":   "XML (RSS)",
	"application/atom+xml":  "XML (Atom)",
	"application/xhtml+xml": "XML (XHTML)",

	// Изображения
	"image/jpeg":    "JPEG изображение",
	"image/png":     "PNG изображение",
	"image/gif":     "GIF изображение",
	"image/tiff":    "TIFF изображение",
	"image/webp":    "WebP изображение",
	"image/bmp":     "BMP изображение",
	"image/svg+xml": "SVG изображение",
}

func validateTenderFile(fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return fmt.Errorf("не удалось определить тип файла: %w", err)
	}

	file.Seek(0, 0)

	if fileHeader.Size == 0 {
		return fmt.Errorf("файл пустой")
	}

	// (50 MB)
	if fileHeader.Size > 50*1024*1024 {
		return fmt.Errorf("файл слишком большой: %d MB (макс. 50 MB)", fileHeader.Size/1024/1024)
	}
	mimeStr := mime.String()
	if _, ok := allowedMimeTypes[mimeStr]; ok {
		return nil
	}
	return fmt.Errorf("недопустимый тип файла: %s (разрешены: PDF, Word, Excel, ZIP, XML, изображения)", mimeStr)
}

func (h *Handlers) CreateTender(IdStatus int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), CreateTenderPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		if IdStatus == 2 {
			f, err := h.readRightsInContext(r.Context(), PublishTenderPerms)
			if err != nil {
				h.handleError(w, "Error get Rights:", err, 500)
				return
			}
			if !f {
				h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
				return
			}
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
			var tender models.Tender
			tender.Name = r.PostForm.Get("name")
			tender.Description = r.PostForm.Get("description")
			var err error
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

			tender.IdCompany = user.IdCompany
			tender.IdStatus = IdStatus //Активный, 1 - Черновик todo в случае с кнопкой сохранить как черновик статус должен быть 1

			//Validation
			err = tender.Validate()
			if err != nil {
				h.handleError(w, "tender validation failed", err, 500)
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

			err = tx.Tenders().Create(r.Context(), &tender)
			if err != nil {
				h.handleError(w, "Failed tender creation", fmt.Errorf("Failed tender creation in transaction: %v", err), 500)
				return
			}
			var log models.Log
			log.IdUser = user.ID
			log.IdEntity = 2
			log.IdType = 1
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
					err = tx.LinkerTCategory(tender.ID).Create(r.Context(), i)
					if err != nil {
						h.handleError(w, "Failed tender-category link creation", fmt.Errorf("Failed tender-category link  creation in transaction: %v", err), 500)
						return
					}
				}
			}

			var path string
			path = fmt.Sprintf("./documents/tender%d", tender.ID)

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
							h.Logger.Errorf("Error of delete tenders mkdir:%s;Error:%v", path, errF)
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
func (h *Handlers) GetCreateTenderWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), CreateTenderPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_create.html")
		if err != nil {
			h.handleError(w, "Failed tendercreate load:", err, 500)
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
		cv := BuildCategoryTree(categories, catlinks)
		dist, err := h.Repo.District().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get districts:", err, 500)
			return
		}
		type Data struct {
			District []models.District
			CatView  []views.CategoryView
		}

		err = tmpl.Execute(w, &Data{District: dist, CatView: cv})
		if err != nil {
			h.handleError(w, "Failed tendercreate render:", err, 500)
			return
		}
	}
}
