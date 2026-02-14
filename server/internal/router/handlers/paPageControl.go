package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetProfilwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
		}
		user, err := h.repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
		}
		company, err := h.repo.Company().GetByID(r.Context(), user.IdCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
		}

		tmpl, err := template.ParseFiles("./client/pages/personal_account.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		type data struct {
			User    models.User
			Company models.Company
		}
		d := &data{User: *user, Company: *company}
		err = tmpl.Execute(w, d)
		if err != nil {
			h.handleError(w, "Failed profil render:", err, 500)
			return
		}
	}
}
func (h *Handlers) EditingLK() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"companyName", "companyEmail", "companyAddress", "companyINN", "companyEGRUL",
				"userLogin", "userName", "userEmail", "userId", "companyId",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}

			cid, err := strconv.Atoi(r.PostForm.Get("companyId"))
			if err != nil {
				h.handleError(w, "bad company id parse", err, 400)
				return
			}
			company, err := h.repo.Company().GetByID(r.Context(), cid)
			if err != nil {
				h.handleError(w, "undefined company", err, 400)
				return
			}

			uid, err := strconv.Atoi(r.PostForm.Get("userId"))
			if err != nil {
				h.handleError(w, "bad user id parse", err, 400)
				return
			}
			user, err := h.repo.Users().GetByID(r.Context(), uid)
			if err != nil {
				h.handleError(w, "undefined user", err, 400)
				return
			}

			company.Name = r.PostForm.Get("companyName")
			company.Email = r.PostForm.Get("companyEmail")
			company.Address = r.PostForm.Get("companyAddress")
			company.INN = r.PostForm.Get("companyINN")
			company.EGRUL = r.PostForm.Get("companyEGRUL")
			company.Description = r.PostForm.Get("companyDescription")
			if err := company.Validate(); err != nil {
				h.handleError(w, "Failed validation company", err, 400)
				return
			}
			tx, err := h.repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			err = tx.Company().Update(r.Context(), company)
			if err != nil {
				h.handleError(w, "Failed update", err, 500)
				return
			}

			user.Login = r.PostForm.Get("userLogin")
			user.Name = r.PostForm.Get("userName")
			user.Email = r.PostForm.Get("userEmail")
			if err := user.Validate(); err != nil {
				h.handleError(w, "Failed validation user", err, 400)
				return
			}
			err = tx.Users().Update(r.Context(), user)
			if err != nil {
				h.handleError(w, "Failed update", err, 500)
				return
			}
			err = tx.Commit()
			if err != nil {
				h.handleError(w, "Failed commit", err, 500)
				return
			}
			w.WriteHeader(204)
			return
		}
		h.handleError(w, "Invalid Content-Type", fmt.Errorf("expected multipart/form-data, got %s", contentType), 400)
		return
	}
}
