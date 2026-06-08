package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Right struct {
	ID      int    `db:"id_right" json:"id_right"`
	Name    string `db:"name" json:"name" validate:"required,min=2,max=50"`
	RusName string `db:"rus_name" json:"rus_name" validate:"required,min=2,max=50"`
	Checked bool
}

func (r *Right) Validate() error {
	err := validate.Struct(r)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
