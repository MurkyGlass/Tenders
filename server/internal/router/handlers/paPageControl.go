package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetProfilwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		role, err := h.Repo.RoleInCompany().GetByID(r.Context(), user.IdRoleInCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/personal_account.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		docs, err := h.Repo.Doc().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get files:", err, 500)
			return
		}
		doclinks, err := h.Repo.LinkerDoc(0).Company().GetAll(r.Context(), company.ID)
		if err != nil {
			h.handleError(w, "Failed get files-links:", err, 500)
			return
		}
		var files []models.Doc
		for _, link := range doclinks {
			for _, d := range docs {
				if d.ID == link.DocID {
					files = append(files, d)
					break
				} else {
					continue
				}
			}
		}
		rights, err := h.Repo.Right().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get rights:", err, 500)
			return
		}
		roles, err := h.Repo.RoleInCompany().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get roles:", err, 500)
			return
		}
		var ownroles []models.RoleInCompany
		for _, rol := range roles {
			if rol.IdCompany != nil {
				if *rol.IdCompany == user.IdCompany {
					ownroles = append(ownroles, rol)
				}
			}
		}
		type data struct {
			User    models.User
			Company models.Company
			Role    models.RoleInCompany
			Files   []models.Doc
			Rights  []models.Right
			Roles   []models.RoleInCompany
		}
		d := &data{User: *user, Company: *company, Role: *role, Files: files, Rights: rights, Roles: ownroles}
		err = tmpl.Execute(w, d)
		if err != nil {
			h.handleError(w, "Failed profil render:", err, 500)
			return
		}
	}
}

// Изминения личных данных пользователя и компании производимое директором или имеющим доступ, для иных случаев следует обдумать отдельный метод
func (h *Handlers) EditingLK() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := h.readRightsInContext(r.Context(), EditCompanyDataPerms)
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
			company, err := h.Repo.Company().GetByID(r.Context(), cid)
			if err != nil {
				h.handleError(w, "undefined company", err, 400)
				return
			}
			idc, ok := r.Context().Value("id_user").(int)
			if !ok {
				h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			}
			curuser, err := h.Repo.Users().GetByID(r.Context(), idc)
			if err != nil {
				h.handleError(w, "db request failed", err, 500)
			}
			uid, err := strconv.Atoi(r.PostForm.Get("userId"))
			if err != nil {
				h.handleError(w, "bad user id parse", err, 400)
				return
			}
			user, err := h.Repo.Users().GetByID(r.Context(), uid)
			if err != nil {
				h.handleError(w, "undefined user", err, 400)
				return
			}
			//
			if curuser.IdCompany != user.IdCompany || curuser.IdCompany != company.ID {
				h.handleError(w, "В доступе отказано", fmt.Errorf("Conflict, id company by user != id company by user"), 409)
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
			id, ok := r.Context().Value("id_user").(int)
			if !ok {
				h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
				return
			}
			err = tx.Log().Create(r.Context(), &models.Log{IdUser: id, IdEntity: 1, IdType: 2}).Company().Create(r.Context(), company.ID)
			if err != nil {
				h.handleError(w, "Failed log company update", err, 500)
				return
			}
			// Действие произвденное над пользователем необязательно совершил данный пользователь!!! ВАЖНО !!! ОБДУМАТЬ!
			_, err = tx.Log().Create(r.Context(), &models.Log{IdUser: id, IdEntity: 5, IdType: 2}).Exists(r.Context())
			if err != nil {
				h.handleError(w, "Failed log user update", err, 500)
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
func (h *Handlers) CreateRoleInCompany() func(w http.ResponseWriter, r *http.Request) {
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
			var role models.RoleInCompany
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
			err = tx.RoleInCompany().Create(r.Context(), &role)
			if err != nil {
				h.handleError(w, "Failed save role", err, 500)
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
func (h *Handlers) CreateNewUser() func(w http.ResponseWriter, r *http.Request) {
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
				"userLogin", "userName", "userEmail", "password", "enterypassword", "role",
			}

			for _, field := range requiredFields {
				if r.PostForm.Get(field) == "" {
					h.handleError(w, "Required field is empty", fmt.Errorf("field %s is required", field), 400)
					return
				}
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
			var user models.User
			user.Email = r.PostForm.Get("userEmail")
			user.Login = r.PostForm.Get("userLogin")
			user.Name = r.PostForm.Get("userName")
			idrol, err := strconv.Atoi(r.PostForm.Get("role"))
			if err != nil {
				h.handleError(w, "Ошибка парсинга должности", fmt.Errorf("Bad role in company"), 500)
				return
			}
			role, err := h.Repo.RoleInCompany().GetByID(r.Context(), idrol)
			if err != nil {
				h.handleError(w, "Должность не найдена", err, 500)
				return
			}
			if role.IdCompany != nil {
				if *role.IdCompany != curuser.IdCompany {
					h.handleError(w, "Должность не относиться к вашей компании", fmt.Errorf("Должность не относиться к вашей компании"), 409)
					return
				}
			} else {
				h.handleError(w, "Системная должность", fmt.Errorf("Должность не относиться к вашей компании"), 409)
				return
			}
			user.IdCompany = curuser.IdCompany
			user.IdRoleInCompany = role.ID
			if r.PostForm.Get("password") != r.PostForm.Get("enterypassword") {
				h.handleError(w, "Пароли не совпадают", fmt.Errorf("Пароли не совпадают"), 500)
				return
			}
			user.Password = r.PostForm.Get("password")
			user.IdRole = 1
			err = user.Validate()
			if err != nil {
				h.handleError(w, "Валидация пользователя провалена", fmt.Errorf("Валидация:%v", err), 500)
				return
			}
			//
			tx, err := h.Repo.BeginTx(r.Context())
			if err != nil {
				h.handleError(w, "Failed begin transaction", err, 500)
				return
			}
			defer tx.Rollback()
			err = tx.Users().Create(r.Context(), &user)
			if err != nil {
				h.handleError(w, "Failed save user", err, 500)
				return
			}
			var log models.Log
			log.DateTime = time.Now()
			log.IdUser = curuser.ID
			log.IdEntity = 5
			log.IdType = 1
			_, err = tx.Log().Create(r.Context(), &log).Exists(r.Context())
			if err != nil {
				h.handleError(w, "Failed Log creation", err, 500)
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
