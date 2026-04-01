package views

import "main/internal/repositories/models"

type TenderView struct {
	ID            int
	Name          string
	Description   string
	DateTimeStart string
	DateTimeEnd   string
	Company       string
	Status        string
	District      string
	Categories    []string
	Files         []models.Doc
}
