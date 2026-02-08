package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Doc struct {
	ID          int    `db:"id_doc" json:"id_doc"`
	Name        string `db:"name" json:"name" validate:"required,min=2,max=50"`
	FileName    string `db:"filename" json:"filename" validate:"required,min=2,max=250"`
	Description string `db:"description" json:"description" validate:"max=500"`
}

func (d *Doc) Validate() error {
	err := validate.Struct(d)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
