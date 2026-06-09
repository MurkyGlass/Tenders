package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetRoleList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), ManageRolesPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}

		tmpl, err := template.ParseFiles("./client/pages/role_list.html")
		if err != nil {
			h.handleError(w, "Failed role list load:", err, 500)
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
		roles, err := h.Repo.RoleInCompany().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get roles", err, 500)
			return
		}
		var viewss []views.RoleInCompanyView
		for _, role := range roles {
			if role.IdCompany != nil {
				if *role.IdCompany == user.IdCompany {
					var view views.RoleInCompanyView
					view.ID = role.ID
					view.Name = role.Name
					links, err := h.Repo.LinkerRoleRight(view.ID).GetById(r.Context())
					if err != nil {
						h.handleError(w, "db request failed", err, 500)
						return
					}
					for _, lin := range links {
						right, err := h.Repo.Right().GetByID(r.Context(), lin)
						if err != nil {
							h.handleError(w, "db request failed", err, 500)
							return
						}
						view.Rights = append(view.Rights, *right)
					}
					viewss = append(viewss, view)
				}
			}
		}
		type Data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Roles            []views.RoleInCompanyView
			Company          *models.Company
		}
		err = tmpl.Execute(w, &Data{LoginForm: LoginForm, RegistrationForm: RegistrationForm, Company: company, Roles: viewss})
		if err != nil {
			h.handleError(w, "Failed role list render:", err, 500)
			return
		}
	}
}
func (h *Handlers) DeleteRole() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), ManageRolesPerms)
		if err != nil {
			h.handleError(w, "Error get Rights:", err, 500)
			return
		}
		if !f {
			h.handleError(w, "No rights for this action", fmt.Errorf("No rights for this action"), 409)
			return
		}
		idrole, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "get id role failed", err, 500)
			return
		}
		role, err := h.Repo.RoleInCompany().GetByID(r.Context(), idrole)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
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
		if user.IdCompany != *role.IdCompany{
			h.handleError(w, "удаление ролей других не предусмотренно", fmt.Errorf("удаление ролей других не предусмотренно"), 409)
			return
		}
		if role.IdCompany == nil || role.ID == 1 || role.ID == 2 {
			h.handleError(w, "удаление системных ролей не предусмотренно", fmt.Errorf("удаление системных ролей не предусмотренно"), 409)
			return
		}
		tx, err := h.Repo.BeginTx(r.Context())
		if err != nil {
			h.handleError(w, "begin tx failed", err, 500)
			return
		}
		defer tx.Rollback()
		err = tx.LinkerRoleRight(role.ID).DeleteAll(r.Context())
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		users, err := tx.Users().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		for _, user := range users {
			if user.IdRoleInCompany == role.ID {
				user.IdRoleInCompany = 2 // no actions
				err = tx.Users().Update(r.Context(), &user)
			}
		}
		err = tx.RoleInCompany().Delete(r.Context(), role.ID)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		err = tx.Commit()
		if err != nil {
			h.handleError(w, "commit failed", err, 500)
			return
		}
		w.WriteHeader(200)
		return
	}
}
