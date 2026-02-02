package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type DistrictRepository interface {
	GetAll(ctx context.Context) ([]models.District, error)
	GetByID(ctx context.Context, id int) (*models.District, error)
}

type districtRepository struct {
	db QueryExecutor
}

func NewDistrictRepository(db QueryExecutor) DistrictRepository {
	return &districtRepository{db: db}
}

func (r *districtRepository) GetAll(ctx context.Context) ([]models.District, error) {
	var districts []models.District
	query := `SELECT * FROM districts`
	err := r.db.SelectContext(ctx, &districts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get districts: %w", err)
	}
	return districts, nil
}
func (r *districtRepository) GetByID(ctx context.Context, id int) (*models.District, error) {
	var d models.District

	query := `SELECT * FROM districts WHERE id_district = $1`

	err := r.db.GetContext(ctx, &d, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get district: %w", err)
	}
	return &d, nil
}
