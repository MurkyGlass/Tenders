package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type DocRepository interface {
	GetAll(ctx context.Context) ([]models.Doc, error)
	GetByID(ctx context.Context, id int) (*models.Doc, error)
	Create(ctx context.Context, doc *models.Doc) error
	Update(ctx context.Context, doc *models.Doc) error
	Delete(ctx context.Context, id int) error
}

type docRepository struct {
	db QueryExecutor
}

func NewDocRepository(db QueryExecutor) DocRepository {
	return &docRepository{db: db}
}

func (r *docRepository) GetAll(ctx context.Context) ([]models.Doc, error) {
	var docs []models.Doc
	query := `SELECT * FROM docs`
	err := r.db.SelectContext(ctx, &docs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get docs: %w", err)
	}
	return docs, nil
}
func (r *docRepository) GetByID(ctx context.Context, id int) (*models.Doc, error) {
	var doc models.Doc

	query := `SELECT * FROM docs WHERE id_doc = $1`

	err := r.db.GetContext(ctx, &doc, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get doc: %w", err)
	}
	return &doc, nil
}
func (r *docRepository) Create(ctx context.Context, doc *models.Doc) error {
	err := doc.Validate()
	if err != nil {
		return err
	}
	//todo add returning id doc
	query := `
		INSERT INTO docs (name, description, filename) 
		VALUES (:name, :description, :filename)
	`
	_, err = r.db.NamedExecContext(ctx, query, doc)
	if err != nil {
		return fmt.Errorf("failed to create doc: %w", err)
	}

	return nil
}
func (r *docRepository) Update(ctx context.Context, doc *models.Doc) error {
	err := doc.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE docs 
		SET name = :name, description = :description, 
			filename = :filename
		WHERE id_doc = :id_doc
	`
	_, err = r.db.NamedExecContext(ctx, query, doc)
	if err != nil {
		return fmt.Errorf("failed to update doc: %w", err)
	}
	return nil
}
func (r *docRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM docs WHERE id_doc = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete doc: %w", err)
	}
	return nil
}
