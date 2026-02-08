package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Category struct {
	ID   int    `db:"id_category" json:"id_category"`
	Name string `db:"name" json:"name" validate:"required,min=2,max=50"`
}

func (c *Category) Validate() error {
	err := validate.Struct(c)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}