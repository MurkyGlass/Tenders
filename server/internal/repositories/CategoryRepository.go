package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Create(ctx context.Context, categ *models.Category, IdParent int) error
	Update(ctx context.Context, name string, Id int, IdParent int) error
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
func (r *categoryRepository) Create(ctx context.Context, categ *models.Category, IdParent int) error {
	err := categ.Validate()
	if err != nil {
		return err
	}
	query := `INSERT INTO Categories (name)
			VALUES (:name)
			RETURNING id_category`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRowxContext(ctx, categ).Scan(&categ.ID)
	if err != nil {
		return fmt.Errorf("failed to create category and get its ID: %w", err)
	}
	if categ.ID == 0 {
		return fmt.Errorf("failed  get category ID: %w", err)
	}
	if IdParent != 0 {
		query := `INSERT INTO Category_Links (id_parent, id_children)
			VALUES ($1,$2)`

		_, err := r.db.ExecContext(ctx, query, IdParent, categ.ID)
		if err != nil {
			return fmt.Errorf("failed to create link by parent: %w", err)
		}
	}
	return nil
}
func (r *categoryRepository) Update(ctx context.Context, name string, Id int, IdParent int) error {
	query := `
		UPDATE Categories
		SET name = $1
		WHERE id_category = $2
	`
	_, err := r.db.ExecContext(ctx, query, name, Id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	if IdParent != -1 {
		query := `Delete from Category_Links where id_children = $1`

		_, err := r.db.ExecContext(ctx, query, Id)
		if err != nil {
			return fmt.Errorf("failed to delete link: %w", err)
		}
		if IdParent != 0 {
			query := `INSERT INTO Category_Links (id_parent, id_children)
			VALUES ($1,$2)`

			_, err := r.db.ExecContext(ctx, query, IdParent, Id)
			if err != nil {
				return fmt.Errorf("failed to create link by parent: %w", err)
			}
		}
	}
	return nil
}
