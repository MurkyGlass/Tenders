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
			h.handleError(w, "Failed main load:", err, 500)
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
			h.handleError(w, "Failed main render:", err, 500)
			return
		}
	}
}
