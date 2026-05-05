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

func (h *Handlers) GetAdminCategoryCreateWindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./client/pages/admin_category_create.html")
		if err != nil {
			h.handleError(w, "Failed admin add category load:", err, 500)
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
			h.handleError(w, "Failed admin add category render:", err, 500)
			return
		}
	}
}

type CategoryRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parentid"` // omitempty позволяет отсутствовать полю
}

func (h *Handlers) AdminCategoryCreate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CategoryRequest
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

		tx, err := h.Repo.BeginTx(r.Context())
		if err != nil {
			h.handleError(w, "db tx failed", err, 500)
			return
		}
		//loging ??
		defer tx.Rollback()
		var cat models.Category
		cat.Name = req.Name
		err = tx.Category().Create(r.Context(), &cat, idparent)
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
