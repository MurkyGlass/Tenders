package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Company struct {
	ID          int    `db:"id_company" json:"id_company"`
	Name        string `db:"name" json:"name" validate:"required,min=3,max=100"`
	Email       string `db:"email" json:"email" validate:"required,email,max=50"`
	Address     string `db:"address" json:"address" validate:"required,min=3,max=100"`
	INN         string `db:"inn" json:"inn" validate:"required,min=10,max=12"`
	EGRUL       string `db:"egrul" json:"egrul" validate:"required,min=13,max=13"`
	Description string `db:"description" json:"description" validate:"max=500"`
}

func (c *Company) Validate() error {
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
