package repositories

import (
	"context"
	"fmt"
)

// Только в TX! Обновление програмно не предусмотренно, следует удалить старую связь и создать новую.
// Get - методы вызывать только через репозиторий!
type LinkerRight interface {
	Create(ctx context.Context, idR int) error
	Delete(ctx context.Context, idT int) error
	Exists(ctx context.Context) (bool, error)
	ExistsByID(ctx context.Context, Id int) (bool, error)
	GetById(ctx context.Context) ([]int, error)
}
type linkerRight struct {
	idRole int
	db     QueryExecutor
}

func NewLinkerRight(idRole int, db QueryExecutor) LinkerRight {
	return &linkerRight{idRole: idRole, db: db}
}

func (l *linkerRight) Create(ctx context.Context, idRight int) error {
	q := `
		INSERT INTO Right_RoleInCompany (id_role,id_right)
		VALUES ($1,$2)
	`
	_, err := l.db.ExecContext(ctx, q, l.idRole, idRight)
	if err != nil {
		return fmt.Errorf("Failed create link role - right: %w", err)
	}
	return nil
}
func (l *linkerRight) Delete(ctx context.Context, idRight int) error {
	q := `
		DELETE FROM Right_RoleInCompany WHERE id_role = $1 AND id_right = $2
	`
	_, err := l.db.ExecContext(ctx, q, l.idRole, idRight)
	if err != nil {
		return fmt.Errorf("Failed delete link role - right: %w", err)
	}
	return nil
}

// return true where this role have links by right, else false. If error, return false and error
func (l *linkerRight) Exists(ctx context.Context) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Right_RoleInCompany WHERE id_role = $1
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idRole).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link role - right: %w", err)
	}
	return ex, nil
}

// return true where this role have link by this right(IdR), else false. If error, return false and error
func (l *linkerRight) ExistsByID(ctx context.Context, IdRight int) (bool, error) {

	q := `
		SELECT EXISTS (
			SELECT 1 FROM Right_RoleInCompany WHERE id_role = $1 AND id_right = $2
		)
	`
	var ex bool
	err := l.db.QueryRowContext(ctx, q, l.idRole, IdRight).Scan(&ex)
	if err != nil {
		return false, fmt.Errorf("Failed find link role - right: %w", err)
	}
	return ex, nil
}
func (l *linkerRight) GetById(ctx context.Context) ([]int, error) {

	q := `
			SELECT id_right FROM Right_RoleInCompany WHERE id_role = $1
	`
	var rigts []int
	err := l.db.SelectContext(ctx, &rigts, q, l.idRole)
	if err != nil {
		return nil, fmt.Errorf("Failed find link role - right: %w", err)
	}
	return rigts, nil
}
