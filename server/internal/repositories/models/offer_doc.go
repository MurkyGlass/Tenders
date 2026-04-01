package models

type Offer_Doc struct {
	OfferID int `json:"offer_id" db:"id_offer"`
	DocID   int `json:"doc_id" db:"id_doc"`
}
