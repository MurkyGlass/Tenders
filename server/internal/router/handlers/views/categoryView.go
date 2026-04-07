package views

import "main/internal/repositories/models"

type CategoryView struct {
	Category models.Category
	Childs   []CategoryView
	Checked  bool
}
