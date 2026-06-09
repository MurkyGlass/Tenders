package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetUserList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), ManageUsersPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}

		tmpl, err := template.ParseFiles("./client/pages/user_list.html")
		if err != nil {
			h.handleError(w, "Failed user list load:", err, 500)
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
		company, err := h.Repo.Company().GetByID(r.Context(), user.IdCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		users, err := h.Repo.Users().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		var usersView []views.UserView
		for _, u := range users {
			if u.IdCompany == user.IdCompany {

				var view views.UserView
				view.ID = u.ID
				view.Name = u.Name
				view.Email = u.Email
				view.Login = u.Login
				role, err := h.Repo.RoleInCompany().GetByID(r.Context(), u.IdRoleInCompany)
				if err != nil {
					h.handleError(w, "db request failed", err, 500)
					return
				}
				view.Role = role.Name
				usersView = append(usersView, view)

			}
		}
		type Data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Users            []views.UserView
			Company          *models.Company
		}
		err = tmpl.Execute(w, &Data{LoginForm: LoginForm, RegistrationForm: RegistrationForm, Company: company, Users: usersView})
		if err != nil {
			h.handleError(w, "Failed user list render:", err, 500)
			return
		}
	}
}
func (h *Handlers) DeleteUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), ManageUsersPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			return
		}
		curuser, err := h.Repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}

		iduser, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "get id user failed", err, 500)
			return
		}
		user, err := h.Repo.Users().GetByID(r.Context(), iduser)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		if user.IdCompany != curuser.IdCompany {
			h.handleError(w, "удаление чужих сотрудников не предусмотренно", fmt.Errorf("удаление сотрудников не предусмотренно"), 409)
			return
		}
		if user.IdRoleInCompany == 1 {
			h.handleError(w, "удаление директора не предусмотренно", fmt.Errorf("удаление директора не предусмотренно"), 409)
			return
		}
		if user.IdRole != 1 {
			h.handleError(w, "удаление администратора или модератора не предусмотренно", fmt.Errorf("удаление администратора или модератора не предусмотренно"), 409)
			return
		}
		tx, err := h.Repo.BeginTx(r.Context())
		if err != nil {
			h.handleError(w, "begin tx failed", err, 500)
			return
		}
		defer tx.Rollback()
		err = tx.Log().DeleteByUser(r.Context(), user.ID)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		err = tx.Users().Delete(r.Context(), user.ID)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}

		err = tx.Commit()
		if err != nil {
			h.handleError(w, "commit failed", err, 500)
			return
		}
		if curuser.ID == user.ID {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
				MaxAge:   -1,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
				MaxAge:   -1,
			})
		}
		w.WriteHeader(200)
		return
	}
}
