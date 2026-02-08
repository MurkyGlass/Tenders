package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
}

type categoryRepository struct {
	db QueryExecutor
}

func NewCategoryRepository(db QueryExecutor) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}
func (r *categoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	var cat models.Category

	query := `SELECT * FROM categories WHERE id_category = $1`

	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &cat, nil
}
