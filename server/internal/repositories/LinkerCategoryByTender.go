package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

// Только в TX! Обновление програмно не предусмотренно, следует удалить старую связь и создать новую.
// Get - методы вызывать только через репозиторий!
type LinkerCategory interface {
	GetAllByTender(ctx context.Context) ([]models.TenderCategory, error)
	GetAll(ctx context.Context) ([]models.TenderCategory, error)
	Create(ctx context.Context, idCat int) error
	Delete(ctx context.Context, idCat int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
}
type linkerCategory struct {
	idTender int
	db       QueryExecutor
}

func NewLinkerCategory(idTender int, db QueryExecutor) LinkerCategory {
	return &linkerCategory{idTender: idTender, db: db}
}
func (l *linkerCategory) GetAll(ctx context.Context) ([]models.TenderCategory, error) {
	q := `
		SELECT * FROM Tender_Category
	`
	var arr []models.TenderCategory
	err := l.db.SelectContext(ctx, &arr, q)
	if err != nil {
		return nil, fmt.Errorf("Failed get tender-category: %w", err)
	}
	return arr, nil
}
func (l *linkerCategory) GetAllByTender(ctx context.Context) ([]models.TenderCategory, error) {
	q := `
		SELECT * FROM Tender_Category WHERE id_tender = $1
	`
	var arr []models.TenderCategory
	err := l.db.SelectContext(ctx, &arr, q, l.idTender)
	if err != nil {
		return nil, fmt.Errorf("Failed get tender-category: %w", err)
	}
	return arr, nil
}
func (l *linkerCategory) Create(ctx context.Context, idCateg int) error {
	q := `
		INSERT INTO Tender_Category (id_tender,id_category)
		VALUES ($1,$2)
	`
	_, err := l.db.ExecContext(ctx, q, l.idTender, idCateg)
	if err != nil {
		return fmt.Errorf("Failed create link tender - category: %w", err)
	}
	return nil
}
func (l *linkerCategory) Delete(ctx context.Context, idCateg int) error {
	q := `
		DELETE FROM Tender_Category WHERE id_tender = $1 AND id_category = $2
	`
	_, err := l.db.ExecContext(ctx, q, l.idTender, idCateg)
	if err != nil {
		return fmt.Errorf("Failed delete link tender - category: %w", err)
	}
	return nil
}

// return true where this tender have links by category, else false. If error, return false and error
func (l *linkerCategory) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Tender_Category WHERE id_tender = $1
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idTender).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link tender - category: %w", err)
	}
	return ex, nil
}

// return true where this tender have link by this category(IdC), else false. If error, return false and error
func (l *linkerCategory) ExistsByID(ctx context.Context, IdC int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Tender_Category WHERE id_tender = $1 AND id_category = $2
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idTender, IdC).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link tender - category: %w", err)
	}
	return ex, nil
}
