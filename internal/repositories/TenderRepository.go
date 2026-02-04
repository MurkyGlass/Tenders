package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type TenderRepository interface {
	GetAll(ctx context.Context) ([]models.Tender, error)
	GetByID(ctx context.Context, id int) (*models.Tender, error)
	Create(ctx context.Context, tender *models.Tender) error
	Update(ctx context.Context, tender *models.Tender) error
	Delete(ctx context.Context, id int) error
}

type tenderRepository struct {
	db QueryExecutor
}

func NewTenderRepository(db QueryExecutor) TenderRepository {
	return &tenderRepository{db: db}
}

func (r *tenderRepository) GetAll(ctx context.Context) ([]models.Tender, error) {
	var tenders []models.Tender
	query := `SELECT * FROM tenders`
	err := r.db.SelectContext(ctx, &tenders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenders: %w", err)
	}
	return tenders, nil
}
func (r *tenderRepository) GetByID(ctx context.Context, id int) (*models.Tender, error) {
	var tender models.Tender

	query := `SELECT * FROM tenders WHERE id_tender = $1`

	err := r.db.GetContext(ctx, &tender, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}
	return &tender, nil
}
func (r *tenderRepository) Create(ctx context.Context, tender *models.Tender) error {
	err := tender.Validate()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO tenders (name, description, datetime_start, datetime_end, id_company, id_status, id_district) 
		VALUES (:name, :description, :datetime_start, :datetime_end, :id_company, :id_status, :id_district)
		RETURNING id_tender
	`
	namedQuery, err := r.db.PrepareNamedContext(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to prepare query(tender): %w", err)
    }
    defer namedQuery.Close()
    
    err = namedQuery.QueryRowContext(ctx, tender).Scan(&tender.ID)
    if err != nil {
        return fmt.Errorf("failed to create tender and get its ID: %w", err)
    }

	return nil
}
func (r *tenderRepository) Update(ctx context.Context, tender *models.Tender) error {
	err := tender.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE tenders 
		SET name = :name, description = :description, 
			datetime_start = :datetime_start, datetime_end = :datetime_end, 
			id_company = :id_company, id_status = :id_status,
			id_district = :id_district
		WHERE id_tender = :id_tender
	`
	_, err = r.db.NamedExecContext(ctx, query, tender)
	if err != nil {
		return fmt.Errorf("failed to update tender: %w", err)
	}
	return nil
}
func (r *tenderRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tenders WHERE id_tender = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tender: %w", err)
	}
	return nil
}
