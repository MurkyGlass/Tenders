package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type StatusRepository interface {
	GetAll(ctx context.Context) ([]models.Status, error)
	GetByID(ctx context.Context, id int) (*models.Status, error)
}

type statusRepository struct {
	db QueryExecutor
}

func NewStatusRepository(db QueryExecutor) StatusRepository {
	return &statusRepository{db: db}
}

func (r *statusRepository) GetAll(ctx context.Context) ([]models.Status, error) {
	var st []models.Status
	query := `SELECT * FROM statuses`
	err := r.db.SelectContext(ctx, &st, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}
	return st, nil
}
func (r *statusRepository) GetByID(ctx context.Context, id int) (*models.Status, error) {
	var st models.Status

	query := `SELECT * FROM statuses WHERE id_status = $1`

	err := r.db.GetContext(ctx, &st, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stasus: %w", err)
	}
	return &st, nil
}
