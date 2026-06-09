package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetUserEdit() func(w http.ResponseWriter, r *http.Request) {
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

		tmpl, err := template.ParseFiles("./client/pages/user_edit.html")
		if err != nil {
			h.handleError(w, "Failed edit user load:", err, 500)
			return
		}

		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			return
		}
		currentuser, err := h.Repo.Users().GetByID(r.Context(), id)
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
		var view views.UserView

		if currentuser.IdCompany == user.IdCompany {
			view.ID = user.ID
			view.Name = user.Name
			view.Email = user.Email
			view.Login = user.Login

			role, err := h.Repo.RoleInCompany().GetByID(r.Context(), user.IdRoleInCompany)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
				return
			}
			if role.IsCreater {
				h.handleError(w, "В доступе отказано", fmt.Errorf("В доступе отказано"), 409)
				return
			}
			view.Role = role.Name

		} else {
			h.handleError(w, "В доступе отказано", fmt.Errorf("В доступе отказано"), 409)
			return
		}
		roles, err := h.Repo.RoleInCompany().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		var comproles []models.RoleInCompany
		for _, rol := range roles {
			if rol.IdCompany != nil {
				if *rol.IdCompany == currentuser.IdCompany {
					comproles = append(comproles, rol)
				}
			} else {
				if rol.ID == 2 {
					comproles = append(comproles, rol)
				}
			}
		}
		type Data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			User             views.UserView
			Roles            []models.RoleInCompany
		}
		err = tmpl.Execute(w, &Data{LoginForm: LoginForm, RegistrationForm: RegistrationForm, Roles: comproles, User: view})
		if err != nil {
			h.handleError(w, "Failed user edit render:", err, 500)
			return
		}
	}
}
func (h *Handlers) EditUser() func(w http.ResponseWriter, r *http.Request) {
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
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"userLogin", "userName", "userEmail", "role",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}
			str := r.PostForm.Get("role")
			idrole, err := strconv.Atoi(str)
			if err != nil {
				h.handleError(w, "Invalid Parsing id role in company by edit user", err, 400)
				return
			}
			//

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
				h.handleError(w, "get id role failed", err, 500)
				return
			}
			user, err := h.Repo.Users().GetByID(r.Context(), iduser)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
				return
			}
			if user.IdCompany != curuser.IdCompany || user.IdRoleInCompany == 1 {
				h.handleError(w, "редактирование пользователей из других компаний или директорских профилей не предусмотренно", fmt.Errorf("редактирование пользователей из других компаний или директорских профилей не предусмотренно"), 409)
				return
			}
			user.Name = r.PostForm.Get("userName")
			user.Email = r.PostForm.Get("userEmail")
			user.Login = r.PostForm.Get("userLogin")
			user.IdRoleInCompany = idrole
			// unique role name in one company
			users, err := h.Repo.Users().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed get users", err, 500)
				return
			}
			for _, ro := range users {
				if (ro.Login == user.Login || ro.Email == user.Email) && ro.ID != user.ID {
					h.handleError(w, "К сожалению такой логин или почта уже существует", fmt.Errorf("Ксожалению такая логин или почта уже существует"), 400)
					return
				}
			}
			//
			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			err = tx.Users().Update(r.Context(), user)
			if err != nil {
				h.handleError(w, "Failed update role", err, 500)
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
