package models

type TenderCategory struct {
	IdTender   int `json:"id_tender" db:"id_tender"`
	IdCategory int `json:"id_category" db:"id_category"`
}
