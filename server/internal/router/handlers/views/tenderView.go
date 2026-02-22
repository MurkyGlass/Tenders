package views

import (
	"time"
)

type TenderView struct {
	ID            int
	Name          string
	Description   string
	DateTimeStart time.Time
	DateTimeEnd   time.Time
	Company       string
	Status        string
	District      string
}
