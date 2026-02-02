package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type RightRepository interface {
	GetAll(ctx context.Context) ([]models.Right, error)
	GetByID(ctx context.Context, id int) (*models.Right, error)
}

type rightRepository struct {
	db QueryExecutor
}

func NewRightRepository(db QueryExecutor) RightRepository {
	return &rightRepository{db: db}
}

func (r *rightRepository) GetAll(ctx context.Context) ([]models.Right, error) {
	var rights []models.Right
	query := `SELECT * FROM rights`
	err := r.db.SelectContext(ctx, &rights, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rights: %w", err)
	}
	return rights, nil
}
func (r *rightRepository) GetByID(ctx context.Context, id int) (*models.Right, error) {
	var right models.Right

	query := `SELECT * FROM rights WHERE id_right = $1`

	err := r.db.GetContext(ctx, &right, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get right: %w", err)
	}
	return &right, nil
}
