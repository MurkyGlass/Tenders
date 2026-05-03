package views

import "main/internal/repositories/models"

type OfferView struct {
	ID             int
	IdTender       int
	Description    string
	Price          float64
	DateTimeCreate string
	Company        string
	Status         string
	IdCompany      int
	Files          []models.Doc
}
