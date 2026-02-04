package repositories

import (
	"context"
	"fmt"
)
//Только в TX!
type Linker interface {
	Company() LCompany
	Tender() LTender
	Offer() LOffer
}
type linkerDoc struct {
	idDoc int
	doc   *docRepository
	err   error
}

type LCompany interface {
	Create(ctx context.Context, idComp int) error
}
type lCompany struct {
	l *linkerDoc
}

func NewLCompany(l *linkerDoc) LCompany {
	return &lCompany{l: l}
}
func (l *linkerDoc) Company() LCompany {
	return NewLCompany(l)
}
func (l *lCompany) Create(ctx context.Context, idComp int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Doc_Company (id_doc,id_company)
		VALUES ($1,$2)
	`
	_, err := l.l.doc.db.ExecContext(ctx, q, l.l.idDoc, idComp)
	if err != nil {
		return fmt.Errorf("Failed create link doc - company: %w", err)
	}
	return nil
}

type LTender interface {
	Create(ctx context.Context, idT int) error
}
type ltender struct {
	l *linkerDoc
}

func NewLTender(l *linkerDoc) LTender {
	return &ltender{l: l}
}
func (l *linkerDoc) Tender() LTender {
	return NewLTender(l)
}
func (l *ltender) Create(ctx context.Context, idT int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Doc_Tender (id_doc,id_tender)
		VALUES ($1,$2)
	`
	_, err := l.l.doc.db.ExecContext(ctx, q, l.l.idDoc, idT)
	if err != nil {
		return fmt.Errorf("Failed create link doc - tender: %w", err)
	}
	return nil
}

type LOffer interface {
	Create(ctx context.Context, idOf int) error
}
type loffer struct {
	l *linkerDoc
}

func NewLOffer(l *linkerDoc) LOffer {
	return &loffer{l: l}
}
func (l *linkerDoc) Offer() LOffer {
	return NewLOffer(l)
}
func (l *loffer) Create(ctx context.Context, idO int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Doc_Offer (id_doc,id_company)
		VALUES ($1,$2)
	`
	_, err := l.l.doc.db.ExecContext(ctx, q, l.l.idDoc, idO)
	if err != nil {
		return fmt.Errorf("Failed create link doc - offer: %w", err)
	}
	return nil
}
