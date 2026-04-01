package models

type Tender_Doc struct {
	TenderID int `json:"tender_id" db:"id_tender"`
	DocID    int `json:"doc_id" db:"id_doc"`
}
