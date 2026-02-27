package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type CategoryLinkRepository interface {
	GetAll(ctx context.Context) ([]models.LinkView, error)
	GetByParentID(ctx context.Context, id int) (*models.LinkView, error)
	GetByChildrenID(ctx context.Context, id int) (*models.LinkView, error)
}

type categoryLinkRepository struct {
	db QueryExecutor
}

func NewCategoryLinkRepository(db QueryExecutor) CategoryLinkRepository {
	return &categoryLinkRepository{db: db}
}

func (r *categoryLinkRepository) GetAll(ctx context.Context) ([]models.LinkView, error) {
	var categories []models.LinkView
	query := `SELECT * FROM Category_Links`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categorieslinks: %w", err)
	}
	return categories, nil
}
func (r *categoryLinkRepository) GetByParentID(ctx context.Context, id int) (*models.LinkView, error) {
	var cat models.LinkView

	query := `SELECT * FROM Category_Links WHERE id_parent = $1`

	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get categorylink: %w", err)
	}
	return &cat, nil
}
func (r *categoryLinkRepository) GetByChildrenID(ctx context.Context, id int) (*models.LinkView, error) {
	var cat models.LinkView

	query := `SELECT * FROM Category_Links WHERE id_children = $1`

	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get categorylink: %w", err)
	}
	return &cat, nil
}