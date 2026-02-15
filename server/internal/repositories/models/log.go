package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Log struct {
	ID       int       `db:"id_log" json:"id_log"`
	IdUser   int       `db:"id_user" json:"id_user" validate:"required"`
	IdEntity int       `db:"id_entity" json:"id_entity" validate:"required"`
	IdType   int       `db:"id_type" json:"id_type" validate:"required"`
	DateTime time.Time `db:"datetime_create" json:"datetime_create"`
}

func (l *Log) Validate() error {
	err := validate.Struct(l)
	if err != nil {
		var valerr []error
		for _, err := range err.(validator.ValidationErrors) {
			valerr = append(valerr, fmt.Errorf("Field %s validation error: %s;", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %w ; %w", err, errors.Join(valerr...))
	}
	return nil
}
