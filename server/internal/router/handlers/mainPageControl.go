package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetMainwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/main.html")
		if err != nil {
			h.handleError(w, "Failed main load:", err, 500)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			h.handleError(w, "Failed main render:", err, 500)
			return
		}
	}
}

func (h *Handlers) Registration() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var company models.Company
		var user models.User
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"companyName", "companyEmail", "companyAddress", "companyINN", "companyEGRUL",
				"userLogin", "userName", "userEmail", "userPassword", "confirmPassword",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}

			if r.PostForm.Get("userPassword") != r.PostForm.Get("confirmPassword") {
				h.handleError(w, "Passwords do not match", fmt.Errorf("passwords do not match"), 400)
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
			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			err = tx.Company().Create(r.Context(), &company)
			if err != nil {
				h.handleError(w, "Failed insert", err, 500)
				return
			}

			user.Login = r.PostForm.Get("userLogin")
			user.Name = r.PostForm.Get("userName")
			user.Email = r.PostForm.Get("userEmail")
			user.Password = r.PostForm.Get("userPassword")
			user.IdRoleInCompany = 1
			user.IdRole = 1
			user.IdCompany = company.ID
			if err := user.Validate(); err != nil {
				h.handleError(w, "Failed validation user", err, 400)
				return
			}
			err = tx.Users().Create(r.Context(), &user)
			if err != nil {
				h.handleError(w, "Failed insert", err, 500)
				return
			}

			err = tx.Log().Create(r.Context(),&models.Log{IdUser: user.ID,IdEntity: 1,IdType: 1}).Company().Create(r.Context(),company.ID)
			if err != nil {
				h.handleError(w, "Failed log company create", err, 500)
				return
			}
			_,err = tx.Log().Create(r.Context(),&models.Log{IdUser: user.ID,IdEntity: 5,IdType: 1}).Exists(r.Context())
			if err != nil {
				h.handleError(w, "Failed log user create", err, 500)
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

