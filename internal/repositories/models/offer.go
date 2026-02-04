package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Offer struct {
	ID             int       `db:"id_offer" json:"id_offer"`
	Description    string    `db:"description" json:"description" validate:"max=500"`
	Price          float64   `db:"price" json:"price" validate:"required"`
	DateTimeCreate time.Time `db:"datetime_create" json:"datetime_create" validate:"required"`
	IdCompany      int       `db:"id_company" json:"id_company" validate:"required"`
	IdStatus       int       `db:"id_status" json:"id_status" validate:"required"`
	IdTender       int       `db:"id_tender" json:"id_tender" validate:"required"`
}

func (o *Offer) Validate() error {
	err := validate.Struct(o)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
