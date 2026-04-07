package repositories

import (
	"context"
	"fmt"
)

// Только в TX! Обновление програмно не предусмотренно, следует удалить старую связь и создать новую.
// Get - методы вызывать только через репозиторий!
type LinkerLog interface {
	Company() LlCompany
	Tender() LlTender
	Offer() LlOffer
	Doc() LlDoc
	Exists(ctx context.Context) (bool, error)
}
type linkerLog struct {
	idLog int
	db    QueryExecutor
	err   error
}

func NewLinkerLog(id int, db QueryExecutor, err error) LinkerLog {
	return &linkerLog{idLog: id, db: db, err: err}
}

// return true where this log have links, else false. If error, return false and error
func (l *linkerLog) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Company WHERE id_log = $1 
			UNION 
			SELECT 1 FROM Log_Tender WHERE id_log = $1
			UNION
			SELECT 1 FROM Log_Offer WHERE id_log = $1
			UNION
			SELECT 1 FROM Log_Doc WHERE id_log = $1
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idLog).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log: %w", err)
	}
	return ex, nil
}

type LlCompany interface {
	Create(ctx context.Context, idComp int) error
	Delete(ctx context.Context, idComp int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type llcompany struct {
	l *linkerLog
}

func NewLlCompany(l *linkerLog) LlCompany {
	return &llcompany{l: l}
}
func (l *linkerLog) Company() LlCompany {
	return NewLlCompany(l)
}
func (l *llcompany) Create(ctx context.Context, idComp int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Log_Company (id_log,id_company)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idComp)
	if err != nil {
		return fmt.Errorf("Failed create link log - company: %w", err)
	}
	return nil
}
func (l *llcompany) Delete(ctx context.Context, idComp int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Log_Company WHERE id_log = $1 AND id_company = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idComp)
	if err != nil {
		return fmt.Errorf("Failed delete link log - company: %w", err)
	}
	return nil
}

// return true where this log have links by company, else false. If error, return false and error
func (l *llcompany) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Company WHERE id_log = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - company: %w", err)
	}
	return ex, nil
}

// return true where this log have link by this company(IdC), else false. If error, return false and error
func (l *llcompany) ExistsByID(ctx context.Context, IdC int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Company WHERE id_log = $1 AND id_company = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog, IdC).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - company: %w", err)
	}
	return ex, nil
}

type LlTender interface {
	Create(ctx context.Context, idT int) error
	Delete(ctx context.Context, idT int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type lltender struct {
	l *linkerLog
}

func NewLlTender(l *linkerLog) LlTender {
	return &lltender{l: l}
}
func (l *linkerLog) Tender() LlTender {
	return NewLlTender(l)
}
func (l *lltender) Create(ctx context.Context, idT int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Log_Tender (id_log,id_tender)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idT)
	if err != nil {
		return fmt.Errorf("Failed create link log - tender: %w", err)
	}
	return nil
}

func (l *lltender) Delete(ctx context.Context, idT int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Log_Tender WHERE id_log = $1 AND id_tender = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idT)
	if err != nil {
		return fmt.Errorf("Failed delete link log - tender: %w", err)
	}
	return nil
}

// return true where this log have links by tender, else false. If error, return false and error
func (l *lltender) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Tender WHERE id_log = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - tender: %w", err)
	}
	return ex, nil
}

// return true where this log have link by this tender(IdT), else false. If error, return false and error
func (l *lltender) ExistsByID(ctx context.Context, IdT int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Tender WHERE id_log = $1 AND id_tender = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog, IdT).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - tender: %w", err)
	}
	return ex, nil
}

type LlOffer interface {
	Create(ctx context.Context, idOf int) error
	Delete(ctx context.Context, idOf int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type lloffer struct {
	l *linkerLog
}

func NewLlOffer(l *linkerLog) LlOffer {
	return &lloffer{l: l}
}
func (l *linkerLog) Offer() LlOffer {
	return NewLlOffer(l)
}
func (l *lloffer) Create(ctx context.Context, idO int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Log_Offer (id_log,id_offer)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idO)
	if err != nil {
		return fmt.Errorf("Failed create link log - offer: %w", err)
	}
	return nil
}
func (l *lloffer) Delete(ctx context.Context, idO int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Log_Offer WHERE id_log = $1 AND id_offer = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idO)
	if err != nil {
		return fmt.Errorf("Failed delete link log - offer: %w", err)
	}
	return nil
}

// return true where this log have links by offer, else false. If error, return false and error
func (l *lloffer) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Offer WHERE id_log = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - offer: %w", err)
	}
	return ex, nil
}

// return true where this log have link by this offer(IdOf), else false. If error, return false and error
func (l *lloffer) ExistsByID(ctx context.Context, IdOf int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Offer WHERE id_log = $1 AND id_offer = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog, IdOf).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - offer: %w", err)
	}
	return ex, nil
}

type LlDoc interface {
	Create(ctx context.Context, idD int) error
	Delete(ctx context.Context, idD int) error
	DeleteAll(ctx context.Context, idD int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type lldoc struct {
	l *linkerLog
}

func NewLlDoc(l *linkerLog) LlDoc {
	return &lldoc{l: l}
}
func (l *linkerLog) Doc() LlDoc {
	return NewLlDoc(l)
}
func (l *lldoc) Create(ctx context.Context, idD int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		INSERT INTO Log_Doc (id_log,id_doc)
		VALUES ($1,$2)
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idD)
	if err != nil {
		return fmt.Errorf("Failed create link log - doc: %w", err)
	}
	return nil
}
func (l *lldoc) Delete(ctx context.Context, idD int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Log_Doc WHERE id_log = $1 AND id_doc = $2
	`
	_, err := l.l.db.ExecContext(ctx, q, l.l.idLog, idD)
	if err != nil {
		return fmt.Errorf("Failed delete link log - doc: %w", err)
	}
	return nil
}
func (l *lldoc) DeleteAll(ctx context.Context, idD int) error {
	if l.l.err != nil {
		return l.l.err
	}
	q := `
		DELETE FROM Log_Doc WHERE id_doc = $1
	`
	_, err := l.l.db.ExecContext(ctx, q, idD)
	if err != nil {
		return fmt.Errorf("Failed delete link log - doc: %w", err)
	}
	return nil
}
// return true where this log have links by doc, else false. If error, return false and error
func (l *lldoc) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Doc WHERE id_log = $1
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - doc: %w", err)
	}
	return ex, nil
}

// return true where this log have link by this doc(IdD), else false. If error, return false and error
func (l *lldoc) ExistsByID(ctx context.Context, IdD int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Log_Doc WHERE id_log = $1 AND id_doc = $2
		)
	`
	var ex bool
	err := l.l.db.QueryRowContext(ctx, q, l.l.idLog, IdD).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link log - doc: %w", err)
	}
	return ex, nil
}
