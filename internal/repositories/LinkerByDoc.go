package repositories

import (
	"context"
	"fmt"
)

// Только в TX! Обновление програмно не предусмотренно, следует удалить старую связь и создать новую.
// Get - методы вызывать только через репозиторий!
type LinkerDoc interface {
	Company() LCompany
	Tender() LTender
	Offer() LOffer
	Exists(ctx context.Context) (bool, error)
}
type linkerDoc struct {
	idDoc int
	db    QueryExecutor
	err   error
}

func NewLinkerDoc(id int, db QueryExecutor, err error) LinkerDoc {
	return &linkerDoc{idDoc: id, db: db, err: err}
}

// return true where this doc have links, else false. If error, return false and error
func (l *linkerDoc) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Company WHERE id_doc = $1 
			UNION 
			SELECT 1 FROM Doc_Tender WHERE id_doc = $1
			UNION
			SELECT 1 FROM Doc_Offer WHERE id_doc = $1
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idDoc).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc: %w", err)
	}
	return ex, nil
}

type LCompany interface {
	Create(ctx context.Context, idComp int) error
	Delete(ctx context.Context, idComp int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type lcompany struct {
	l *linkerDoc
}

func NewLCompany(l *linkerDoc) LCompany {
	return &lcompany{l: l}
}
func (l *linkerDoc) Company() LCompany {
	return NewLCompany(l)
}
func (l *lcompany) Create(ctx context.Context, idComp int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Doc_Company (id_doc,id_company)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idComp)
	if err != nil {
		return fmt.Errorf("Failed create link doc - company: %w", err)
	}
	return nil
}
func (l *lcompany) Delete(ctx context.Context, idComp int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Doc_Company WHERE id_doc = $1 AND id_company = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idComp)
	if err != nil {
		return fmt.Errorf("Failed delete link doc - company: %w", err)
	}
	return nil
}

// return true where this doc have links by company, else false. If error, return false and error
func (l *lcompany) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Company WHERE id_doc = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - company: %w", err)
	}
	return ex, nil
}

// return true where this doc have link by this company(IdC), else false. If error, return false and error
func (l *lcompany) ExistsByID(ctx context.Context, IdC int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Company WHERE id_doc = $1 AND id_company = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc, IdC).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - company: %w", err)
	}
	return ex, nil
}

type LTender interface {
	Create(ctx context.Context, idT int) error
	Delete(ctx context.Context, idT int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
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
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idT)
	if err != nil {
		return fmt.Errorf("Failed create link doc - tender: %w", err)
	}
	return nil
}

func (l *ltender) Delete(ctx context.Context, idT int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Doc_Tender WHERE id_doc = $1 AND id_tender = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idT)
	if err != nil {
		return fmt.Errorf("Failed delete link doc - tender: %w", err)
	}
	return nil
}

// return true where this doc have links by tender, else false. If error, return false and error
func (l *ltender) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Tender WHERE id_doc = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - tender: %w", err)
	}
	return ex, nil
}

// return true where this doc have link by this tender(IdT), else false. If error, return false and error
func (l *ltender) ExistsByID(ctx context.Context, IdT int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Tender WHERE id_doc = $1 AND id_tender = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc, IdT).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - tender: %w", err)
	}
	return ex, nil
}

type LOffer interface {
	Create(ctx context.Context, idOf int) error
	Delete(ctx context.Context, idOf int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
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
		INSERT INTO Doc_Offer (id_doc,id_offer)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idO)
	if err != nil {
		return fmt.Errorf("Failed create link doc - offer: %w", err)
	}
	return nil
}
func (l *loffer) Delete(ctx context.Context, idO int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Doc_Offer WHERE id_doc = $1 AND id_offer = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idDoc, idO)
	if err != nil {
		return fmt.Errorf("Failed delete link doc - offer: %w", err)
	}
	return nil
}

// return true where this doc have links by offer, else false. If error, return false and error
func (l *loffer) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Offer WHERE id_doc = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - offer: %w", err)
	}
	return ex, nil
}

// return true where this doc have link by this offer(IdOf), else false. If error, return false and error
func (l *loffer) ExistsByID(ctx context.Context, IdOf int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Doc_Offer WHERE id_doc = $1 AND id_offer = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idDoc, IdOf).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link doc - offer: %w", err)
	}
	return ex, nil
}
