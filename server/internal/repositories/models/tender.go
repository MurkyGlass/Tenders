package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Tender struct {
	ID            int       `db:"id_tender" json:"id_tender"`
	Name          string    `db:"name" json:"name" validate:"required,min=3,max=100"`
	Description   string    `db:"description" json:"description" validate:"max=500"`
	DateTimeStart time.Time `db:"datetime_start" json:"datetime_start" validate:"required"`
	DateTimeEnd   time.Time `db:"datetime_end" json:"datetime_end" validate:"required"`
	IdCompany     int       `db:"id_company" json:"id_company" validate:"required"`
	IdStatus      int       `db:"id_status" json:"id_status" validate:"required"`
	IdDistrict    int       `db:"id_district" json:"id_district" validate:"required"`
}

func (t *Tender) Validate() error {
	err := validate.Struct(t)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
