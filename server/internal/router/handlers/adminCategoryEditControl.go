package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func (h *Handlers) GetAdminCategoryEditWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/admin_category_edit.html")
		if err != nil {
			h.handleError(w, "Failed admin update category load:", err, 500)
			return
		}
		categories, err := h.Repo.Category().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		type Data struct {
			Categories []models.Category
		}
		err = tmpl.Execute(w, &Data{Categories: categories})
		if err != nil {
			h.handleError(w, "Failed admin update category render:", err, 500)
			return
		}
	}
}

type CategoryUpRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentid"`
}

func (h *Handlers) AdminCategoryEdit() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CategoryUpRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		idparent, err := strconv.Atoi(req.ParentID)
		if err != nil {
			h.handleError(w, "Parent id conv failed", fmt.Errorf("Parent id conv failed"), 500)
			return
		}
		idcateg, err := strconv.Atoi(req.ID)
		if err != nil {
			h.handleError(w, "Category id conv failed", fmt.Errorf("Category id conv failed"), 500)
			return
		}
		category, err := h.Repo.Category().GetByID(r.Context(), idcateg)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		category.Name = req.Name
		err = category.Validate()
		if err != nil {
			h.handleError(w, "validate failed", err, 500)
			return
		}

		tx, err := h.Repo.BeginTx(r.Context())
		if err != nil {
			h.handleError(w, "db tx failed", err, 500)
			return
		}
		//loging ??
		defer tx.Rollback()

		err = tx.Category().Update(r.Context(), category.Name, category.ID, idparent)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		err = tx.Commit()
		if err != nil {
			h.handleError(w, "db request commit failed", err, 500)
			return
		}
		w.WriteHeader(201)
	}
}
