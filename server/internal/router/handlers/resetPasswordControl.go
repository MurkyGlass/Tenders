package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/email"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (h *Handlers) GetResetWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/reset_password.html")
		if err != nil {
			h.handleError(w, "Failed reset load:", err, 500)
			return
		}
		type Data struct {
		}
		err = tmpl.Execute(w, &Data{})
		if err != nil {
			h.handleError(w, "Failed reset render:", err, 500)
			return
		}
	}
}
func (h *Handlers) GetResetForm() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/reset_password_form.html")
		if err != nil {
			h.handleError(w, "Failed reset-form load:", err, 500)
			return
		}
		token := r.URL.Query().Get("token")
		reset, err := h.Repo.Reset().GetByToken(r.Context(), token)
		if err != nil {
			h.handleError(w, "Failed get Token", err, 409)
			return
		}
		if reset == nil {
			h.handleError(w, "Token is invalid or expires", err, 409)
			return
		}

		type Data struct {
		}
		err = tmpl.Execute(w, &Data{})
		if err != nil {
			h.handleError(w, "Failed reset-form render:", err, 500)
			return
		}
	}
}
func (h *Handlers) UpdatePassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		token := r.URL.Query().Get("token")
		reset, err := h.Repo.Reset().GetByToken(r.Context(), token)
		if err != nil {
			h.handleError(w, "Token is invalid or expires", err, 500)
			return
		}

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"password", "doublepassword",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}
			pass := r.PostForm.Get("password")
			dpass := r.PostForm.Get("doublepassword")
			if pass != dpass {
				h.handleError(w, "Passwords is not identical", err, 400)
				return
			}
			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			user, err := tx.Users().GetByID(r.Context(), reset.UserID)
			if err != nil {
				h.handleError(w, "Failed get user by reset token", err, 500)
				return
			}

			user.Password = pass
			err = tx.Users().UpdatePassword(r.Context(), user)
			if err != nil {
				h.handleError(w, "Failed update password ", err, 500)
				return
			}
			err = tx.Reset().DeleteByToken(r.Context(), reset.Token)
			if err != nil {
				h.handleError(w, "Failed delete token", err, 500)
				return
			}
			err = tx.Commit()
			if err != nil {
				h.handleError(w, "Failed commit ", err, 500)
				return
			}
			w.WriteHeader(204)
			return
		}
		h.handleError(w, "Invalid Content-Type", fmt.Errorf("expected multipart/form-data, got %s", contentType), 400)
		return
	}
}
func (h *Handlers) ResetPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"mail",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}

			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			mail := r.PostForm.Get("mail")
			user, err := tx.Users().GetByEmail(r.Context(), mail)
			if err != nil {
				h.handleError(w, "Failed get user by email", err, 500)
				return
			}
			if user == nil {
				h.handleError(w, "User not found", err, 404)
				return
			}

			var reset models.ResetToken
			reset.Token = uuid.New().String()
			reset.UserID = user.ID
			reset.ExpiresAt = time.Now().Add(15 * time.Minute)

			err = tx.Reset().Create(r.Context(), &reset)
			if err != nil {
				h.handleError(w, "Failed insert ", err, 500)
				return
			}
			err = email.SendResetEmail(user.Email, reset.Token)
			if err != nil {
				h.handleError(w, "Failed send mail ", err, 500)
				return
			}
			err = tx.Commit()
			if err != nil {
				h.handleError(w, "Failed commit ", err, 500)
				return
			}
			w.WriteHeader(204)
			return
		}
		h.handleError(w, "Invalid Content-Type", fmt.Errorf("expected multipart/form-data, got %s", contentType), 400)
		return
	}
}
