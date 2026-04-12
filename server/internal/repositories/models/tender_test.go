package models_test

import (
	"main/internal/repositories/models"
	"testing"
	"time"
)

func TestTenderValidation(t *testing.T) {
    now := time.Now()
    future := now.Add(24 * time.Hour)

    tests := []struct {
        name    string
        tender  models.Tender
        wantErr bool
    }{
        {
            name: "Valid tender",
            tender: models.Tender{
                Name:          "Поставка оборудования",
                Description:   "Краткое описание",
                DateTimeStart: now,
                DateTimeEnd:   future,
                IdCompany:     1,
                IdStatus:      2,
                IdDistrict:    1,
            },
            wantErr: false,
        },
        {
            name: "Valid tender(Not Description)",
            tender: models.Tender{
                Name:          "Поставка оборудования",
                Description:   "",
                DateTimeStart: now,
                DateTimeEnd:   future,
                IdCompany:     1,
                IdStatus:      2,
                IdDistrict:    1,
            },
            wantErr: false,
        },
        {
            name: "Empty name",
            tender: models.Tender{
                Name:          "",
                Description:   "Описание",
                DateTimeStart: now,
                DateTimeEnd:   future,
                IdCompany:     1,
                IdStatus:      2,
                IdDistrict:    1,
            },
            wantErr: true,
        },
        {
            name: "Name too long (101 chars)",
            tender: models.Tender{
                Name:          string(make([]byte, 101)),
                Description:   "Описание",
                DateTimeStart: now,
                DateTimeEnd:   future,
                IdCompany:     1,
                IdStatus:      2,
                IdDistrict:    1,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.tender.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}