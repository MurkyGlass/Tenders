package handler

import (
	"fmt"
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetRoleEdit() func(w http.ResponseWriter, r *http.Request) {
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

		tmpl, err := template.ParseFiles("./client/pages/role_edit.html")
		if err != nil {
			h.handleError(w, "Failed role in company edit load:", err, 500)
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
		var view views.RoleInCompanyView
		if role.IdCompany != nil {
			if *role.IdCompany == user.IdCompany {
				view.ID = role.ID
				view.Name = role.Name
				links, err := h.Repo.LinkerRoleRight(view.ID).GetById(r.Context())
				if err != nil {
					h.handleError(w, "db request failed", err, 500)
					return
				}
				rights, err := h.Repo.Right().GetAll(r.Context())
				if err != nil {
					h.handleError(w, "db request failed", err, 500)
					return
				}
				for _, right := range rights {
					right.Checked = false
					for _, lin := range links {
						if right.ID == lin {
							right.Checked = true
						}
					}
					view.Rights = append(view.Rights, right)
				}

			} else {
				h.handleError(w, "В доступе отказано", fmt.Errorf("В доступе отказано"), 409)
				return
			}
		} else {
			h.handleError(w, "В доступе отказано", fmt.Errorf("В доступе отказано"), 409)
			return
		}

		type Data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Role             views.RoleInCompanyView
		}
		err = tmpl.Execute(w, &Data{LoginForm: LoginForm, RegistrationForm: RegistrationForm, Role: view})
		if err != nil {
			h.handleError(w, "Failed role in company edit render:", err, 500)
			return
		}
	}
}
func (h *Handlers) EditRoleInCompany() func(w http.ResponseWriter, r *http.Request) {
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
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}

			// Проверка обязательных полей
			requiredFields := []string{
				"name",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
			}
			IdsStr := r.PostForm.Get("rights")
			if IdsStr == "" {
				h.handleError(w, "Права не назначены", fmt.Errorf("Rights is empety"), 400)
				return
			}
			IDsStrArr := r.PostForm["rights"]
			var Ids []int
			for _, str := range IDsStrArr {
				id, err := strconv.Atoi(str)
				if err != nil {
					h.handleError(w, "Invalid Parsing category id by Filter", err, 400)
					return
				}
				Ids = append(Ids, id)
			}

			//

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
			if role.IdCompany == nil || role.ID == 1 || role.ID == 2 {
				h.handleError(w, "редактирование системных ролей не предусмотренно", fmt.Errorf("редактирование системных ролей не предусмотренно"), 409)
				return
			}
			role.Name = r.PostForm.Get("name")
			role.IsCreater = false
			idc := user.IdCompany
			role.IdCompany = &idc
			err = role.Validate()
			if err != nil {
				h.handleError(w, "ошибка валидации должности", err, 400)
				return
			}
			// unique role name in one company
			roles, err := h.Repo.RoleInCompany().GetAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed get roles", err, 500)
				return
			}
			for _, ro := range roles {
				if ro.IdCompany != nil {
					if *ro.IdCompany == *role.IdCompany && ro.Name == role.Name {
						h.handleError(w, "К сожалению такая должность уже существует", fmt.Errorf("Ксожалению такая должность уже существует"), 400)
						return
					}
				}
			}
			//
			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			err = tx.RoleInCompany().Update(r.Context(), role)
			if err != nil {
				h.handleError(w, "Failed update role", err, 500)
				return
			}
			err = tx.LinkerRoleRight(role.ID).DeleteAll(r.Context())
			if err != nil {
				h.handleError(w, "Failed delete rights for role", err, 500)
				return
			}
			for _, i := range Ids {
				err = tx.LinkerRoleRight(role.ID).Create(r.Context(), i)
				if err != nil {
					h.handleError(w, "Failed save right for role", err, 500)
					return
				}
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
