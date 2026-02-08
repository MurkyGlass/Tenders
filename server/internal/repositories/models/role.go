package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Role struct {
	ID   int    `db:"id_role"`
	Name string `db:"name" validate:"required,min=2,max=50"`
}

func (r *Role) Validate() error {
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
