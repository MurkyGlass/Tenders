package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID              int    `db:"id_user" json:"id_user"`
	Login           string `db:"login" json:"login" validate:"required,min=4,max=30"`
	Name            string `db:"name" json:"name" validate:"required,min=2,max=30"`
	Email           string `db:"email" json:"email" validate:"required,email,max=50"`
	IdRoleInCompany int    `db:"id_role_in_company" json:"id_role_in_company" validate:"required"`
	IdCompany       int    `db:"id_company" json:"id_company" validate:"required"`
	IdRole          int    `db:"id_role" json:"id_role" validate:"required"`
	Password        string `db:"password" json:"password" validate:"required,min=8"`
}

func (u *User) Validate() error {
	err := validate.Struct(u)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
