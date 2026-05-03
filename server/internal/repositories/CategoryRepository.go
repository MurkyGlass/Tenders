package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Create(ctx context.Context, name string, IdParent int) error
	Update(ctx context.Context, name string, Id int) error
	Delete(ctx context.Context, id int) error
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
func (r *categoryRepository) Create(ctx context.Context, name string, IdParent int) error {
	query := `INSERT INTO Categories (name)
			VALUES ($1)
			RETURNING id_category`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	ID := 0
	err = stmt.QueryRowxContext(ctx, name).Scan(&ID)
	if err != nil {
		return fmt.Errorf("failed to create category and get its ID: %w", err)
	}
	if ID == 0 {
		return fmt.Errorf("failed  get category ID: %w", err)
	}
	if IdParent != 0 {
		query := `INSERT INTO Category_Links (id_parent, id_children)
			VALUES ($1,$2)`

		_, err := r.db.ExecContext(ctx, query, IdParent, ID)
		if err != nil {
			return fmt.Errorf("failed to create link by parent: %w", err)
		}
	}
	return nil
}
func (r *categoryRepository) Update(ctx context.Context, name string, Id int) error {
	query := `
		UPDATE Categories
		SET name = $1
		WHERE id_category = $2
	`
	_, err := r.db.ExecContext(ctx, query, name, Id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}
func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Category_Links WHERE id_parent = $1 or id_children = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete catlinks: %w", err)
	}
	query = `DELETE FROM Categories WHERE id_category = $1`
	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
