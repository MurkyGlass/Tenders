package models

type Company_Doc struct {
	CompanyID int `json:"company_id" db:"id_company"`
	DocID     int `json:"doc_id" db:"id_doc"`
}
