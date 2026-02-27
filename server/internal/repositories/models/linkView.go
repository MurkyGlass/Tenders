package models

type LinkView struct {
	FirstID  int `json:"first_id" db:"id_parent"`
	SecondID int `json:"second_id" db:"id_children"`
}
