package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type RoleInCompany struct {
	ID        int    `db:"id_role" json:"id_role"`
	Name      string `db:"name" json:"name" validate:"required,min=2,max=50"`
	IdCompany *int   `db:"id_company" json:"id_company"`
	IsCreater bool   `db:"is_creater"`
}

func (r *RoleInCompany) Validate() error {
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
