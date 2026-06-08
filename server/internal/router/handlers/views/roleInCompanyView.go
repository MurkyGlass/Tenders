package views

import (
	"main/internal/repositories/models"
)

type RoleInCompanyView struct {
	ID        int
	Name      string
	Rights    []models.Right
}
